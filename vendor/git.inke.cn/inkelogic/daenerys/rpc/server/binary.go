package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.inke.cn/BackendPlatform/golang/logging"
	"golang.org/x/net/context"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/ikiosocket"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/metadata"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/rpcerror"
	"git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/tpc/inf/go-tls"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"git.inke.cn/tpc/inf/metrics"
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openlog "github.com/opentracing/opentracing-go/log"
)

type binaryServer struct {
	opts     Options
	router   *router
	server   *ikiosocket.Server
	plugins  []Plugin
	cfg      *config.Register
	shutdown int32
	stop     chan struct{}
	once     sync.Once
}

func BinaryServer(options ...Option) Server {
	opts := newOptions(options...)
	h := &binaryServer{}
	h.router = newRouter()
	h.opts = opts
	h.server = ikiosocket.NewServer(opts.Kit.G(), h.serveBinary)
	h.shutdown = 0
	h.stop = make(chan struct{})
	return h
}

func (h *binaryServer) serveBinary(remoteAddr string, request *ikiosocket.Context) (*ikiosocket.Context, error) {
	var codec = h.opts.Codec
	var start = time.Now()

	optname := "RPC Server"
	elog := log.WithPrefix(
		h.opts.Error,
		"component", "rpcserver-binary",
		"time", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"remote", remoteAddr,
	)
	ctx := tracing.BinaryToContext(h.opts.Tracer, request.Header, optname, elog)
	ns := namespace.GetNamespace(ctx)
	if ns != "" {
		elog = log.WithPrefix(elog, "namespace", ns)
	}

	tls.SetContext(ctx)
	peer := metric.GetSDName(ctx)

	span := opentracing.SpanFromContext(ctx)
	span.SetTag("proto", "rpc/binary")
	defer span.Finish()
	defer tls.Flush()

	makeFailedResponse := func(code int, id uint64, desc, name string) *ikiosocket.Context {
		ext.Error.Set(span, true)
		meta := &metadata.RpcMeta{
			Type:       metadata.RpcMeta_RESPONSE.Enum(),
			SequenceId: proto.Uint64(id),
			Failed:     proto.Bool(true),
			ErrorCode:  proto.Int32(int32(code)),
			Reason:     proto.String(desc),
		}
		body, _ := proto.Marshal(meta)
		if ns != "" {
			metrics.Timer(name, start, metrics.TagCode, code, "namespace", ns)
		} else {
			metrics.Timer(name, start, metrics.TagCode, code)
		}

		return &ikiosocket.Context{
			Header: map[string]string{
				metadata.MetaHeaderKey: string(body),
			},
		}
	}

	if _, ok := request.Header[metadata.MetaHeaderKey]; !ok {
		elog.Log("err", "mate key not exit")
		return makeFailedResponse(rpcerror.Internal, 0, "Meta Key Not Exist", "SServer"), nil
	}
	if _, ok := request.Header[metadata.DataHeaderKey]; !ok {
		elog.Log("err", "data key not exit")
		return makeFailedResponse(rpcerror.Internal, 0, "Data Key Not Exist", "SServer"), nil
	}

	span.LogFields(openlog.String("event", "rpcmate"))
	meta := metadata.RpcMeta{}
	if err := proto.Unmarshal([]byte(request.Header[metadata.MetaHeaderKey]), &meta); err != nil {
		elog.Log("err", err)
		span.LogFields(openlog.Error(err))
		return makeFailedResponse(rpcerror.Internal, 0, "Malfored Metadata", "SServer"), nil
	}
	names := meta.GetMethod()
	pos := strings.LastIndex(names, ".")
	if pos == -1 {
		elog.Log("name", names, "err", "Malfored Method Name")
		return makeFailedResponse(rpcerror.Internal, meta.GetSequenceId(), "Malfored Method Name", "SServer"), nil
	}

	service := names[:pos]
	method := names[pos+1:]
	metricname := fmt.Sprintf("SServer.%s.%s", service, method)
	span.SetOperationName(fmt.Sprintf("RPC Server %s.%s", service, method))

	stype, mtype, args, err := h.router.signature(service, method)
	if err != nil {
		elog.Log("err", err)
		span.LogFields(openlog.Error(err))
		return makeFailedResponse(rpcerror.Internal, meta.GetSequenceId(), fmt.Sprintf("Unknown Service %v", service), metricname), nil
	}
	body := []byte(request.Header[metadata.DataHeaderKey])

	c := core.New()
	rpcctx := &Context{
		core:       c,
		opts:       h.opts,
		Ctx:        ctx,
		Service:    service,
		Method:     method,
		Header:     make(map[string]string),
		RemoteAddr: remoteAddr,
		Body:       body,
		Request:    args,
		Code:       int32(rpcerror.Success),
		Namespace:  ns,
		Peer:       peer,
	}

	blog := log.WithPrefix(
		h.opts.Kit.B(),
		"time", log.DefaultTimestamp,
		"end", "server",
		"trace", log.TraceID(ctx),
		"remote", remoteAddr,
		"cost", log.Cost(start),
		"service", fmt.Sprintf("%s.%s", service, method),
		"proto", "binary",
		"peer", peer,
	)

	if ns != "" {
		blog = log.WithPrefix(blog, "namespace", ns)
	}

	span.LogFields(openlog.String("event", "decode"))
	if err := codec.Decode(body, rpcctx.Request); err != nil {
		blog.Log("err", err, "code", rpcerror.Internal)
		span.LogFields(openlog.Error(err))
		return makeFailedResponse(rpcerror.Internal, meta.GetSequenceId(), fmt.Sprintf("Parse Error %v", err), metricname), nil
	}

	for _, plugin := range h.plugins {
		plugin := plugin

		c.Use(core.Function(func(ctx context.Context, c core.Core) {
			plugin(rpcctx)
		}))
	}

	c.Use(core.Function(func(ctx context.Context, c core.Core) {
		response, err := h.router.call(rpcctx.Ctx, stype, mtype, rpcctx.Request)
		if err != nil {
			c.AbortErr(err) // TODO
			return
		}
		rpcctx.Response = response
	}))

	span.LogFields(openlog.String("event", "call"))
	rpcctx.Next()
	if err := rpcctx.Err(); err != nil {
		blog.Log("err", err, "code", rpcerror.FromUser)
		span.LogFields(openlog.Error(err))
		return makeFailedResponse(rpcerror.FromUser, meta.GetSequenceId(), err.Error(), metricname), nil
	}

	successResponse := func(header map[string]string, body []byte) *ikiosocket.Context {
		meta := &metadata.RpcMeta{
			Type:       metadata.RpcMeta_RESPONSE.Enum(),
			SequenceId: proto.Uint64(meta.GetSequenceId()),
			Failed:     proto.Bool(false),
			ErrorCode:  proto.Int32(int32(rpcerror.Success)),
		}
		m, err := proto.Marshal(meta)
		if err != nil {
			panic(err)
		}
		header[metadata.MetaHeaderKey] = string(m)
		header[metadata.DataHeaderKey] = string(body)
		return &ikiosocket.Context{
			Header: header,
		}
	}

	span.LogFields(openlog.String("event", "decode"))
	body, err = codec.Encode(rpcctx.Response)
	if err != nil {
		blog.Log("err", err, "code", rpcerror.Internal)
		span.LogFields(openlog.Error(err))
		return makeFailedResponse(rpcerror.Internal, meta.GetSequenceId(), err.Error(), metricname), nil
	}

	// reload code,maybe changed by outside logic
	atomic.LoadInt32(&rpcctx.Code)
	blog.Log("err", "<nil>", "code", rpcctx.Code)

	if ns != "" {
		metrics.Timer(metricname, start, metrics.TagCode, rpcctx.Code, "namespace", ns)
	} else {
		metrics.Timer(metricname, start, metrics.TagCode, rpcctx.Code)
	}
	return successResponse(rpcctx.Header, body), nil
}

func (h *binaryServer) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	return h.router.NewHandler(handler, opts...)
}

func (h *binaryServer) Handle(handler Handler) error {
	return h.router.Handle(handler)
}

func (h *binaryServer) Use(list ...Plugin) Server {
	h.plugins = append(h.plugins, list...)
	return h
}

func (h *binaryServer) Start() error {
	var err error
	h.once.Do(func() {
		ln, e := net.Listen("tcp4", h.opts.Address)
		if e != nil {
			err = e
			return
		}

		addr := strings.Split(ln.Addr().String(), ":")
		port, _ := strconv.Atoi(addr[1])
		cfg, e := utils.Register(
			h.opts.Manager, h.opts.Name, "rpc", h.opts.Tags, config.LocalIPString(), port)
		if e != nil {
			err = e
			return
		}

		h.cfg = cfg
		err = h.server.Start(ln)
		if strings.Contains(err.Error(), "use of closed network connection") {
			err = nil
		}
		logging.Infof("waiting for rpc server stop done")
		fmt.Println("waiting for rpc server stop done")
		// waiting for stop done
		<-h.stop
	})
	return err
}

func (h *binaryServer) Stop() error {
	if !atomic.CompareAndSwapInt32(&h.shutdown, 0, 1) {
		return nil
	}

	defer close(h.stop)

	if m := h.opts.Manager; m != nil {
		m.Deregister()
	}
	h.server.Stop()
	return nil
}
