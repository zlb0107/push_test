package client

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/breaker"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"git.inke.cn/inkelogic/daenerys/internal/kit/ratelimit"
	"git.inke.cn/inkelogic/daenerys/internal/kit/recovery"
	"git.inke.cn/inkelogic/daenerys/internal/kit/retry"
	"git.inke.cn/inkelogic/daenerys/internal/kit/sd"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"time"
)

type generalClient struct {
	defaultCore core.Core
	opts        Options
	endpoint    string
}

func newGeneralClient(f sd.Factory, endpoint string, opts Options) *generalClient {
	defaultCore := core.New(
		// recovery
		recovery.Recovery(log.Stdout()),

		// tracing
		tracing.TraceClient(opts.Tracer, fmt.Sprintf("RPC Client %s", endpoint)),

		// peer name
		metric.SDName(opts.SDName),

		// metric
		metric.Metric(fmt.Sprintf("client.%s", endpoint)),

		// slow log
		core.Function(func(ctx context.Context, c core.Core) {
			start := time.Now()
			c.Next(ctx)
			cost := time.Since(start)
			if cost <= opts.Slow {
				return
			}
			span := opentracing.SpanFromContext(ctx)
			span.SetTag("slow", true)
			rpcctx := ctx.Value(rpcContextKey).(*rpcContext)
			log.WithPrefix(
				opts.Kit.S(),
				"time", log.DefaultTimestamp,
				"trace", log.TraceID(ctx),
				"client", opts.Name,
				"cost", log.Cost(start),
			).Log("method", rpcctx.Endpoint, "host", rpcctx.host)
		}),

		// bussiness log
		core.Function(func(ctx context.Context, c core.Core) {
			span := opentracing.SpanFromContext(ctx)
			span.SetOperationName(fmt.Sprintf("RPC Client %s", endpoint))
			start := time.Now()
			c.Next(ctx)
			rpcctx := ctx.Value(rpcContextKey).(*rpcContext)
			log.WithPrefix(
				opts.Kit.B(),
				"time", log.DefaultTimestamp,
				"end", "client",
				"trace", log.TraceID(ctx),
				"name", opts.Name,
				"cost", log.Cost(start),
			).Log("method", rpcctx.Endpoint, "remote", rpcctx.host, "err", c.Err())
		}),

		// rate limter
		ratelimit.Limiter(endpoint, opts.Ratelimit),

		// breaker
		breaker.Breaker(endpoint, opts.Breaker),

		// retry
		retry.Retry(opts.Retries),

		// namespace
		namespace.Namespace(opts.Namespace),

		// service discovery
		sd.Upstream(f, opts.Cluster),
	)

	return &generalClient{
		opts:        opts,
		endpoint:    endpoint,
		defaultCore: defaultCore,
	}
}

func (r *generalClient) Invoke(ctx context.Context, request interface{}, response interface{}, opts ...CallOption) error {
	rpcctx := &rpcContext{
		Endpoint: r.endpoint,
		Request:  request,
		Response: response,
	}
	c := r.defaultCore.Copy()
	c.Next(context.WithValue(ctx, rpcContextKey, rpcctx))
	return c.Err()
}
