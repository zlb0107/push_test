package namespace

import (
	"git.inke.cn/inkelogic/daenerys/internal/core"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

const (
	NAMESPACE = "_namespace_appkey_"
)

func Namespace(namespace string) core.Plugin {
	return core.Function(func(ctx context.Context, c core.Core) {
		if span := opentracing.SpanFromContext(ctx); span != nil && namespace != "" {
			span.SetBaggageItem(NAMESPACE, namespace)
		}
		c.Next(ctx)
	})
}

func GetNamespace(ctx context.Context) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		return span.BaggageItem(NAMESPACE)
	}
	return ""
}
