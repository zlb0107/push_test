package client

import (
	"net/http"
	"time"

	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
	"github.com/opentracing/opentracing-go"
)

type Options struct {
	logger              log.Kit
	tracer              opentracing.Tracer
	dialTimeout         time.Duration
	idleConnTimeout     time.Duration
	keepAliveTimeout    time.Duration
	keepAlivesDisable   bool
	requestTimeout      time.Duration
	retryTimes          int
	maxIdleConnsPerHost int
	maxIdleConns        int
	cluster             *upstream.Cluster
	client              *http.Client
	ratelimit           *ratelimit.Config
	breaker             *breaker.Config
	namespace           string
}

type Option func(*Options)

func newOptions(options ...Option) Options {
	v := Options{}
	for _, o := range options {
		o(&v)
	}

	if v.logger == nil {
		v.logger = log.NewKit(
			log.Stdout(), // bus
			log.Stdout(), // gen
			log.Stdout(), // acc
			log.Stdout(), // slow
		)
	}
	if v.tracer == nil {
		v.tracer = opentracing.GlobalTracer()
	}

	if v.dialTimeout == 0 {
		v.dialTimeout = defaultDialTimeout
	}
	if v.keepAliveTimeout == 0 {
		v.keepAliveTimeout = defaultKeepAliveTimeout
	}
	if v.requestTimeout == 0 {
		v.requestTimeout = defaultRequestTimeout
	}
	if v.idleConnTimeout == 0 {
		v.idleConnTimeout = defaultIdleConnTimeout
	}
	if v.maxIdleConns == 0 {
		v.maxIdleConns = defaultMaxIdleConns
	}
	if v.maxIdleConnsPerHost == 0 {
		v.maxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	}
	return v
}

func Breaker(config *breaker.Config) Option {
	return func(o *Options) {
		o.breaker = config
	}
}

func Ratelimit(config *ratelimit.Config) Option {
	return func(o *Options) {
		o.ratelimit = config
	}
}

func Logger(logger log.Kit) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func Tracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		if tracer != nil {
			o.tracer = tracer
		}
	}
}

func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = d
	}
}

func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.requestTimeout = d
	}
}

func IdleConnTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.idleConnTimeout = d
	}
}

func KeepAliveTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.keepAliveTimeout = d
	}
}

// if true, do not re-use of TCP connections
func KeepAlivesDisable(t bool) Option {
	return func(o *Options) {
		o.keepAlivesDisable = t
	}
}

func RetryTimes(d int) Option {
	return func(o *Options) {
		o.retryTimes = d
	}
}

func MaxIdleConnsPerHost(d int) Option {
	return func(o *Options) {
		o.maxIdleConnsPerHost = d
	}
}

func MaxIdleConns(d int) Option {
	return func(o *Options) {
		o.maxIdleConns = d
	}
}

func Cluster(clu *upstream.Cluster) Option {
	return func(o *Options) {
		if clu != nil {
			o.cluster = clu
		}
	}
}

func WithClient(client *http.Client) Option {
	return func(o *Options) {
		if client != nil {
			o.client = client
		}
	}
}

func Namespace(n string) Option {
	return func(o *Options) {
		o.namespace = n
	}
}
