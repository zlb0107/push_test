package client

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/ikiosocket"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/rpcerror"
	"github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
	openlog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

func SClient(endpoint string, options ...Option) Client {
	opts := newOptions(options...)
	return newGeneralClient(&binaryFactory{
		opts: opts,
		pool: newPool(opts.PoolSize, opts.dialer, opts.PoolTTL, nil),
	}, endpoint, opts)
}

type binaryFactory struct {
	opts Options
	pool *pool
}

func (r *binaryFactory) Name() string {
	return "binary"
}

func (r *binaryFactory) Factory(host string) (core.Plugin, error) {
	sock, err := r.pool.getSocket(host)
	if err != nil {
		return nil, err
	}

	return core.Function(func(ctx context.Context, c core.Core) {
		var (
			err    error
			span   = opentracing.SpanFromContext(ctx)
			codec  = r.opts.Codec
			rpcctx = ctx.Value(rpcContextKey).(*rpcContext)
		)
		rpcctx.host = host
		// TODO
		elog := log.WithPrefix(
			r.opts.Error,
			"component", "binary-client",
			"time", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"trace", log.TraceID(ctx),
		)
		span.SetTag("proto", "rpc/binary")

		defer func() {
			if err == nil {
				return
			}
			ext.Error.Set(span, true)
			c.AbortErr(err)
		}()

		span.LogFields(openlog.String("event", "encode"))
		body, err := codec.Encode(rpcctx.Request)
		if err != nil {
			span.LogFields(
				openlog.String("event", "decode error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, fmt.Errorf("encode: %v", err))
			return
		}

		carrier := opentracing.TextMapCarrier{}
		if tracer := r.opts.Tracer; tracer != nil {
			if err := tracer.Inject(span.Context(), opentracing.TextMap, carrier); err != nil {
				elog.Log("err", err)
			}
		}

		header := make(map[string]string)
		carrier.ForeachKey(func(key, value string) error {
			header[key] = value
			return nil
		})

		span.LogFields(openlog.String("event", "transport"))

		body, err = sock.Call(rpcctx.Endpoint, header, body)
		if err == ikiosocket.ErrExited {
			// only close socket when exit
			r.pool.release(host, sock, err)
		}

		if err != nil {
			elog.Log("err", err)
			span.LogFields(
				openlog.String("event", "transport error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, err)
			return
		}

		span.LogFields(openlog.String("event", "decode"))
		err = codec.Decode(body, rpcctx.Response)
		if err != nil {
			span.LogFields(
				openlog.String("event", "decode error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, fmt.Errorf("decode: %v", err))
			return
		}
	}), nil
}
