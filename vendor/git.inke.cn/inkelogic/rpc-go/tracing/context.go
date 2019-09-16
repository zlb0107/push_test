package tracing

import (
	tls "git.inke.cn/tpc/inf/go-tls"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

func CaptureTraceContext(c context.Context) context.Context {
	if c == nil {
		c = context.Background()
	}
	if span := opentracing.SpanFromContext(c); span == nil {
		if ctx, ok := tls.GetContext(); ok {
			span = opentracing.SpanFromContext(ctx)
			return opentracing.ContextWithSpan(c, span)
		}
	}
	return c
}
