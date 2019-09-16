package metric

import (
	"git.inke.cn/inkelogic/daenerys/internal/core"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

// SDName
func SDName(name string) core.Plugin {
	return core.Function(func(ctx context.Context, c core.Core) {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span.SetBaggageItem("peer_discovery_name", name)
		}
		c.Next(ctx)
	})
}

func GetSDName(ctx context.Context) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		s := span.BaggageItem("peer_discovery_name")
		if s == "" {
			return "unknown"
		}
		return s
	}
	return "unknown"
}
