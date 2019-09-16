package server

import (
	"time"

	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	"github.com/opentracing/opentracing-go"
)

const (
	HTTPReadTimeout  = 60 * time.Second
	HTTPWriteTimeout = 60 * time.Second
	HTTPIdleTimeout  = 90 * time.Second
)

type Options struct {
	logger log.Kit
	tracer opentracing.Tracer

	serviceName string
	port        int

	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration //server keep conn

	certFile string
	keyFile  string

	tags map[string]string

	manager   *registry.ServiceManager
	registry  registry.Backend
	breaker   *breaker.Config
	ratelimit *ratelimit.Config
}

type Option func(*Options)

func newOptions(options ...Option) Options {
	opts := Options{}
	for _, o := range options {
		o(&opts)
	}

	if opts.logger == nil {
		opts.logger = log.NewKit(
			log.Stdout(), //bus
			log.Stdout(), //gen
			log.Stdout(), //acc
			log.Stdout(), //slow
		)
	}

	if opts.tracer == nil {
		opts.tracer = opentracing.GlobalTracer()
	}

	if opts.readTimeout == 0 {
		opts.readTimeout = HTTPReadTimeout
	}
	if opts.writeTimeout == 0 {
		opts.writeTimeout = HTTPWriteTimeout
	}
	if opts.idleTimeout == 0 {
		opts.idleTimeout = HTTPIdleTimeout
	}
	return opts
}

func Ratelimit(config *ratelimit.Config) Option {
	return func(o *Options) {
		o.ratelimit = config
	}
}

func Breaker(config *breaker.Config) Option {
	return func(o *Options) {
		o.breaker = config
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

func Port(port int) Option {
	return func(o *Options) {
		o.port = port
	}
}

func Name(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

// 从连接被接受(accept)到request body完全被读取(如果你不读取body，那么时间截止到读完header为止)
// 包括了TCP消耗的时间,读header时间
// 对于 https请求，ReadTimeout 包括了TLS握手的时间
func ReadTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = d
	}
}

// 从request header的读取结束开始，到response write结束为止 (也就是 ServeHTTP 方法的声明周期)
func WriteTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = d
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = d
	}
}

func CertFile(file string) Option {
	return func(o *Options) {
		o.certFile = file
	}
}

func KeyFile(file string) Option {
	return func(o *Options) {
		o.keyFile = file
	}
}

func Tags(tags map[string]string) Option {
	return func(o *Options) {
		o.tags = tags
	}
}

func Manager(re *registry.ServiceManager) Option {
	return func(o *Options) {
		o.manager = re
	}
}

func Registry(r registry.Backend) Option {
	return func(o *Options) {
		if r != nil {
			o.registry = r
		}
	}
}
