package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/breaker"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"git.inke.cn/inkelogic/daenerys/internal/kit/ratelimit"
	"git.inke.cn/inkelogic/daenerys/internal/kit/recovery"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/tpc/inf/go-upstream/circuit"
	"git.inke.cn/tpc/inf/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context"
)

// core plugin encapsulation
type HandlerFunc func(c *Context)

func (p HandlerFunc) Do(ctx context.Context, flow core.Core) {
	c := ctx.Value(httpServerInternalContext).(*Context)
	c.Ctx = ctx // original ctx maybe changed
	p(c)
}

func ratelimitPlugin(c *Context) {
	if c.opts.ratelimit == nil {
		c.Next()
		return
	}
	limit := c.opts.ratelimit.LimterWithPeer(c.Path, c.Peer)
	if limit != nil && !limit.Allow() {
		c.AbortErr(ratelimit.ErrLimited)
		return
	}
	c.Next()
}

func (s *server) recover() core.Plugin {
	return recovery.Recovery(log.Stdout())
}

const internalReqBodyLogTag = "req_body"
const internalRespBodyLogTag = "resp_body"

func (s *server) logging() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		start := time.Now()
		flow.Next(ctx)
		alog := log.WithPrefix(
			s.options.logger.A(),
			"time", log.DefaultTimestamp,
			"pid", log.PID(),
			"caller", log.DefaultCaller,
			"cost", log.Cost(start),
			"trace_id", log.TraceID(ctx),
			"type", "httpserver",
		)

		cc := ctx.Value(httpServerInternalContext).(*Context)
		logItems := []interface{}{
			"uri", cc.Request.URL.String(),
			"http_code", cc.Response.Status(),
			"req_method", cc.Request.Method,
			"real_ip", getRemoteIP(cc.Request),
			"busi_code", cc.BusiCode(),
			"error", flow.Err(),
		}

		// request body
		if _, ok := cc.loggingExtra[internalReqBodyLogTag]; !ok {
			body := fmt.Sprintf("%q", cc.bodyBuff.Bytes())
			logItems = append(logItems, "req_body", body)
		}
		// response body
		if _, ok := cc.loggingExtra[internalRespBodyLogTag]; !ok {
			logItems = append(logItems, "resp_body", cc.Response.StringBody())
		}

		if len(cc.loggingExtra) > 0 {
			extraList := make([]interface{}, 0)
			for k, v := range cc.loggingExtra {
				extraList = append(extraList, k, v)
			}
			if len(extraList) > 0 {
				logItems = append(logItems, extraList...)
			}
		}
		// will be sorted by json encoder
		_ = alog.Log(logItems...)
	})
}

func (s *server) metric() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		start := time.Now()
		flow.Next(ctx)
		cc := ctx.Value(httpServerInternalContext).(*Context)
		err := flow.Err()
		switch err.(type) {
		case circuit.BreakerError:
			cc.SetBusiCode(BreakerOpen)
		}

		code := cc.BusiCode()
		methodName := strings.Trim(strings.Replace(cc.Request.URL.Path, "/", ".", -1), ".")
		metricPrefix := fmt.Sprintf("RestServe.%s", methodName)
		metricPeerPrefix := fmt.Sprintf("RestServePeer.%s", methodName)
		metrics.Timer(metricPeerPrefix, start, metrics.TagCode, code, "peer", cc.Peer)
		if ns := namespace.GetNamespace(ctx); ns != "" {
			metrics.Timer(metricPrefix, start, metrics.TagCode, code, "namespace", ns)
		} else {
			metrics.Timer(metricPrefix, start, metrics.TagCode, code)
		}
		cc.LoggingExtra("method_name", methodName, "peer", cc.Peer)
	})
}

func (s *server) tracing() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		flow.Next(ctx)
		span := opentracing.SpanFromContext(ctx)
		cc := ctx.Value(httpServerInternalContext).(*Context)
		ext.HTTPStatusCode.Set(span, uint16(cc.Response.Status()))
		ip, _, _ := net.SplitHostPort(cc.Request.RemoteAddr)
		ext.PeerHostIPv4.SetString(span, ip)
		span.SetTag("inkelogic.code", cc.BusiCode())
		cc.LoggingExtra("req_remote_ip", ip)
	})
}

func (s *server) breaker(path string) core.Plugin {
	return breaker.BreakerServer(path, s.options.breaker)
}

func (s *server) flow(ctx *Context) core.Core {
	t := s.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != ctx.Request.Method {
			continue
		}
		root := t[i].root
		//plugin, urlparam, found, matchPath expression
		flow, params, _, mpath := root.getValue(ctx.Request.URL.Path, ctx.Params, false)
		if flow != nil {
			ctx.Params = params
			ctx.Path = mpath
			return flow
		}
		break
	}
	return nil
}

func (s *server) methodNotAllowed(ctx *Context) bool {
	// 405
	t := s.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == ctx.Request.Method {
			continue
		}
		root := t[i].root
		// plugin, urlparam, found, matchPath expression
		flow, _, _, _ := root.getValue(ctx.Request.URL.Path, ctx.Params, false)
		if flow != nil {
			return true
		}
	}
	return false
}
