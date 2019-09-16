package server

import (
	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	"git.inke.cn/inkelogic/daenerys/rpc/codec"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	opentracing "github.com/opentracing/opentracing-go"
)

type Option func(*Options)

type Options struct {
	Codec     codec.Codec
	Address   string
	Name      string
	Tracer    opentracing.Tracer
	Kit       log.Kit
	Manager   *registry.ServiceManager
	Tags      map[string]string
	Error     log.Logger
	Registry  registry.Backend
	Ratelimit *ratelimit.Config
	Breaker   *breaker.Config
}

func newOptions(opt ...Option) Options {
	opts := Options{
		Error: log.Stdout(),
	}
	for _, o := range opt {
		o(&opts)
	}

	if len(opts.Address) == 0 {
		opts.Address = "127.0.0.1:10000"
	}

	if len(opts.Name) == 0 {
		opts.Name = "daenerys-client"
	}

	if opts.Tracer == nil {
		opts.Tracer = opentracing.NoopTracer{}
	}

	if opts.Kit == nil {
		opts.Kit = log.NewKit(
			log.Stdout(),
			log.Stdout(),
			log.Stdout(),
			log.Stdout(),
		)
	}
	return opts
}

func Breaker(config *breaker.Config) Option {
	return func(o *Options) {
		o.Breaker = config
	}
}

func Ratelimit(config *ratelimit.Config) Option {
	return func(o *Options) {
		o.Ratelimit = config
	}
}

// TODO
func Registry(r registry.Backend) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// TODO
func Error(e log.Logger) Option {
	return func(o *Options) {
		o.Error = e
	}
}

// TODO
func Tags(tags map[string]string) Option {
	return func(o *Options) {
		o.Tags = tags
	}
}

// TODO
func Manager(b *registry.ServiceManager) Option {
	return func(o *Options) {
		o.Manager = b
	}
}

// Logger
func LoggerKit(kit log.Kit) Option {
	return func(o *Options) {
		o.Kit = kit
	}
}

// Tracer
func Tracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		o.Tracer = tracer
	}
}

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Codec to use to encode/decode requests for a given content type
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}
