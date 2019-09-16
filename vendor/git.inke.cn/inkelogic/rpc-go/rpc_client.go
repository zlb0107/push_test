package rpc

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"git.inke.cn/inkelogic/daenerys"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/rpc/client"
	"git.inke.cn/inkelogic/daenerys/rpc/codec"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"github.com/opentracing/opentracing-go"
)

type ClientConn struct {
	factory   client.Factory
	endpoints sync.Map
	mu        sync.Mutex
}

var defaultmaxConnsPerHost = 50

var ErrClientConnClosing = errors.New("rpc: the client connection is closing")
var ErrClientLB = errors.New("rpc: config loadbalance not found ip")

type dialOptions struct {
	maxIdleConns  int
	configBalance bool
}

type DialOption func(*dialOptions)

func WithMaxIdleConns(max int) DialOption {
	return func(o *dialOptions) {
		o.maxIdleConns = max
	}
}

func WithConfigBalance(flag bool) DialOption {
	return func(o *dialOptions) {
		o.configBalance = flag
	}
}

func DialService(service string, config RequestOptionIntercace) (*ClientConn, error) {
	return &ClientConn{factory: daenerys.RPCFactory(context.TODO(), service)}, nil
}

func DialServiceContent(ctx context.Context, config RequestOptionIntercace, sc ServerClient) (conn *ClientConn, err error) {
	return &ClientConn{factory: daenerys.RPCFactory(context.TODO(), sc.ServiceName)}, nil
}

func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}

func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	dopts := dialOptions{}
	dopts.maxIdleConns = defaultmaxConnsPerHost
	for _, opt := range opts {
		opt(&dopts)
	}
	cconfig := config.NewCluster()
	cconfig.Name = target
	cconfig.StaticEndpoints = target
	cconfig.Proto = "rpc"
	if err := daenerys.Default.Clusters.InitService(cconfig); err != nil {
		return nil, err
	}

	c := client.HFactory(
		client.Cluster(GetCluster(target)),
		client.Tracer(opentracing.GlobalTracer()),
		client.RequestTimeout(time.Millisecond*100),
		client.Slow(time.Millisecond*50),
		client.Codec(codec.NewJSONCodec()),
		client.Name(target),
		client.Kit(log.NewKit(
			log.Default(filepath.Join(logdirGlobal, "bussiness.log")).Logger(),
			log.Default(filepath.Join(logdirGlobal, "gen.log")).Logger(),
			log.Default(filepath.Join(logdirGlobal, "access.log")).Logger(),
			log.Default(filepath.Join(logdirGlobal, "slow.log")).Logger(),
		)),
	)
	return &ClientConn{factory: c}, nil
}
