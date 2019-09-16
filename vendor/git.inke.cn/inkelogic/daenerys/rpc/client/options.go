package client

import (
	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	"git.inke.cn/inkelogic/daenerys/rpc/codec"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
	"github.com/opentracing/opentracing-go"
	"time"
)

type Options struct {
	Retries   int
	Kit       log.Kit
	Error     log.Logger
	Tracer    opentracing.Tracer
	Codec     codec.Codec
	Name      string
	SDName    string
	Slow      time.Duration
	Ratelimit *ratelimit.Config
	Breaker   *breaker.Config
	Cluster   *upstream.Cluster
	Namespace string

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration
	// Transport Dial Timeout
	DialTimeout time.Duration

	// Default CallOptions
	CallOptions CallOptions

	// only for testing
	dialer dialer

	maxIdleConnsPerHost int
	maxIdleConns        int
	keepAlivesDisable   bool
}

type CallOptions struct {
	// Number of Call attempts
	Retries int
	// Request/Response timeout
	RequestTimeout time.Duration
}

func newOptions(options ...Option) Options {
	opts := Options{
		Retries:     1,
		Kit:         log.NewKit(log.Stdout(), log.Stdout(), log.Stdout(), log.Stdout()),
		Error:       log.Stdout(),
		Tracer:      opentracing.NoopTracer{},
		PoolSize:    DefaultPoolSize,
		PoolTTL:     DefaultPoolTTL,
		Codec:       codec.NewProtoCodec(),
		DialTimeout: DefaultDialTimeout,
		Slow:        time.Millisecond * 30,
		CallOptions: CallOptions{
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
		},
	}
	opts.dialer = defaultDialer{opts}

	for _, o := range options {
		o(&opts)
	}

	if opts.maxIdleConns == 0 {
		opts.maxIdleConns = 50
	}

	if opts.maxIdleConnsPerHost == 0 {
		opts.maxIdleConnsPerHost = 50
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

func Slow(slow time.Duration) Option {
	return func(o *Options) {
		o.Slow = slow
	}
}

func Error(elog log.Logger) Option {
	return func(o *Options) {
		o.Error = elog
	}
}

func SDName(n string) Option {
	return func(o *Options) {
		o.SDName = n
	}
}

func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Logger sets the logger
func Kit(logger log.Kit) Option {
	return func(o *Options) {
		o.Kit = logger
	}
}

func Cluster(cluster *upstream.Cluster) Option {
	return func(o *Options) {
		o.Cluster = cluster
	}
}

// Codec sets the codec
func Codec(codec codec.Codec) Option {
	return func(o *Options) {
		o.Codec = codec
	}
}

// Tracer sets the opentracing tracer
func Tracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		o.Tracer = tracer
	}
}

// PoolSize sets the connection pool size
func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

// PoolSize sets the connection pool size
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

// Transport dial timeout
func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = d
	}
}

// Number of retries when making the request.
// Should this be a Call Option?
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// The request timeout.
// Should this be a Call Option?
func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}

// WithRetries is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

// WithRequestTimeout is a CallOption which overrides that which
// set in Options.CallOptions
func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
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

func KeepAlivesDisable(t bool) Option {
	return func(o *Options) {
		o.keepAlivesDisable = t
	}
}
