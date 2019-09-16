package tracing

import (
	"git.inke.cn/inkelogic/daenerys/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context"
)

func BinaryToContext(tracer opentracing.Tracer, header map[string]string, operationName string, logger log.Logger) context.Context {
	var span opentracing.Span
	wireContext, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(header))
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		logger.Log("err", err)
	}
	span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	return opentracing.ContextWithSpan(context.Background(), span)
}
