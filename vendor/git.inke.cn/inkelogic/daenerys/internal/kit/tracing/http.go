package tracing

import (
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"

	"git.inke.cn/inkelogic/daenerys/log"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context"
)

// ContextToHTTP returns an http RequestFunc that injects an OpenTracing Span
// found in `ctx` into the http headers. If no such Span can be found, the
// RequestFunc is a noop.
func ContextToHTTP(ctx context.Context, tracer opentracing.Tracer, req *http.Request) error {
	// Try to find a Span in the Context.
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	//http trace
	r := &requestTracer{sp: span}
	nCtx := httptrace.WithClientTrace(ctx, r.clientTrace())
	req = req.WithContext(nCtx)

	// Add standard OpenTracing tags.
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	host, portString, err := net.SplitHostPort(req.URL.Host)
	if err == nil {
		ext.PeerHostname.Set(span, host)
		if port, err := strconv.Atoi(portString); err != nil {
			ext.PeerPort.Set(span, uint16(port))
		}
	} else {
		ext.PeerHostname.Set(span, req.URL.Host)
	}

	// There's nothing we can do with any errors here.
	if err = tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	); err != nil {
		return err
	}

	return nil
}

// HTTPToContext returns an http RequestFunc that tries to join with an
// OpenTracing trace found in `req` and starts a new Span called
// `operationName` accordingly. If no trace could be found in `req`, the Span
// will be a trace root. The Span is incorporated in the returned Context and
// can be retrieved with opentracing.SpanFromContext(ctx).
func HTTPToContext(tracer opentracing.Tracer, req *http.Request, operationName string, logger log.Logger) context.Context {
	// Try to join to a trace propagated in `req`.
	var span opentracing.Span
	wireContext, err := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		// TODO
		logger.Log("err", err)
	}

	span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	return opentracing.ContextWithSpan(context.Background(), span)
}
