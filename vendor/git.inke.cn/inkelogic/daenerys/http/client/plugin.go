package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"net/http/httputil"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/breaker"
	"git.inke.cn/inkelogic/daenerys/internal/kit/ratelimit"
	"git.inke.cn/inkelogic/daenerys/internal/kit/recovery"
	"git.inke.cn/inkelogic/daenerys/internal/kit/retry"
	"git.inke.cn/inkelogic/daenerys/internal/kit/sd"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

// core plugin encapsulation
type HandlerFunc func(c *Context)

func (p HandlerFunc) Do(ctx context.Context, flow core.Core) {
	c := ctx.Value(httpClientInternalContext).(*Context)
	c.Ctx = ctx // original ctx maybe changed
	p(c)
}

func (c *Context) recover() core.Plugin {
	return recovery.Recovery(log.Stdout())
}

func (c *Context) tracing(operation string) core.Plugin {
	return tracing.TraceClient(c.client.options.tracer, operation)
}

func (c *Context) retry() core.Plugin {
	if c.client.options.retryTimes > 0 {
		return retry.Retry(c.client.options.retryTimes)
	}
	return nil
}

func (c *Context) upstream() core.Plugin {
	if c.client.options.cluster != nil {
		return sd.Upstream(c, c.client.options.cluster)
	}
	return nil
}

func (c *Context) breaker(service, path string) core.Plugin {
	_ = service
	return breaker.Breaker(path, c.client.options.breaker)
}

func (c *Context) ratelimit(path string) core.Plugin {
	return ratelimit.Limiter(path, c.client.options.ratelimit)
}

const businessLogger = "_internal_business_log_"

func (c *Context) logging() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		start := time.Now()
		blog := log.WithPrefix(
			c.client.options.logger.B(),
			"time", log.DefaultTimestamp,
			"pid", log.PID(),
			"caller", log.DefaultCaller,
			"cost", log.Cost(start),
			"trace_id", log.TraceID(ctx),
		)
		ctx = context.WithValue(ctx, businessLogger, blog)
		flow.Next(ctx)
		endTime := time.Now()
		cc := ctx.Value(httpClientInternalContext).(*Context)
		var rspCode int
		if cc.Resp != nil {
			rspCode = cc.Resp.Code()
		} else if flow.Err() != nil {
			rspCode = http.StatusInternalServerError
		}

		method := cc.Req.raw.Method
		path := cc.Req.raw.URL.Path
		address := cc.Req.raw.URL.Host
		serviceName := cc.Req.ro.serviceName
		operation := fmt.Sprintf("HTTP Client %s %s", method, path)
		span := opentracing.SpanFromContext(ctx)
		span.SetOperationName(operation)
		span.SetTag("http_code", rspCode)
		logItems := []interface{}{
			"type", "httpclient",
			"service", serviceName,
			"req_method", method,
			"req_path", path,
			"http_code", rspCode,
			"address", address,
			"err", flow.Err(),
		}
		_ = blog.Log(logItems...)

		// slow log
		costTime := endTime.Sub(start)
		if cc.Req.ro != nil && cc.Req.ro.slowTime > 0 && costTime > time.Duration(cc.Req.ro.slowTime)*time.Millisecond {
			span.SetTag("slow", true)
			slog := log.WithPrefix(
				c.client.options.logger.S(),
				"time", log.DefaultTimestamp,
				"pid", log.PID(),
				"caller", log.DefaultCaller,
				"cost", log.Cost(start),
				"trace_id", log.TraceID(ctx),
			)
			_ = slog.Log(logItems...)
		}
	})
}

func (c *Context) sender() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		span := opentracing.SpanFromContext(ctx)
		span.LogFields(opentracinglog.String("event", "HttpRequestSending"))
		span.SetTag("client_retry", c.client.options.retryTimes)
		doRetry := 0
		needRetry := 0
		blog := ctx.Value(businessLogger).(log.Logger)
		cc := ctx.Value(httpClientInternalContext).(*Context)
		if cc.Req.ro != nil && cc.Req.ro.retryTimes > 0 {
			needRetry = cc.Req.ro.retryTimes
			span.SetTag("request_retry", needRetry)
		}

		timeout := 0
		if cc.Req.ro != nil && cc.Req.ro.reqTimeout > 0 {
			timeout = cc.Req.ro.reqTimeout
			span.SetTag("request_timeout", timeout)
		}

		// record request
		reqBuf, _ := httputil.DumpRequest(cc.Req.raw, true)

		// request body
		var bodyBytes []byte
		if cc.Req.raw.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(cc.Req.raw.Body)
		}

		nCtx := ctx
		var cancel context.CancelFunc

	RETRY: // todo:this retry way should change http host
		cc.Req.WithBody(bytes.NewReader(bodyBytes))
		if timeout > 0 {
			nCtx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		}
		cc.Req.raw = cc.Req.raw.WithContext(nCtx)

		// http client do
		resp, err := c.client.client.Do(cc.Req.raw)
		if err != nil || (resp != nil && resp.StatusCode >= http.StatusInternalServerError) {
			blog.Log("httpdo_error", err)
			if needRetry > 0 && doRetry < needRetry {
				span.LogFields(
					opentracinglog.String("event", "retry"),
					opentracinglog.String("times", fmt.Sprintf("%d/%d", doRetry, needRetry)),
					opentracinglog.Error(err))
				doRetry++

				// drop
				if resp != nil {
					io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
				}

				// cancel nCtx
				if cancel != nil {
					cancel()
				}
				goto RETRY
			}
		}

		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				span.LogFields(
					opentracinglog.String("event", "httpdo"),
					opentracinglog.String("reason", "Context DeadlineExceeded"),
					opentracinglog.Error(err))
				if needRetry > 0 { //needRetry为特定request重试逻辑,此处忽略外层client重试逻辑
					b := retry.BreakError{
						Err: context.DeadlineExceeded,
					}
					flow.AbortErr(b)
				} else { // 使用外部client重试逻辑
					flow.AbortErr(context.DeadlineExceeded)
				}
			} else if strings.Contains(err.Error(), "context canceled") {
				span.LogFields(
					opentracinglog.String("event", "httpdo"),
					opentracinglog.String("reason", "Context Canceled"),
					opentracinglog.Error(err))
				flow.AbortErr(context.Canceled)
			} else {
				span.LogFields(
					opentracinglog.String("event", "httpdo"),
					opentracinglog.String("reason", fmt.Sprintf("url.Error wraped error: %s", reflect.TypeOf(err).String())),
					opentracinglog.Error(err))
				err = &url.Error{
					Err: err,
				}
				flow.AbortErr(err)
			}
			ext.Error.Set(span, true)
		}

		respBuf := utils.DumpRespBody(resp)
		cc.Resp, err = BuildResp(cc.Req.raw, resp)
		if err != nil {
			blog.Log("buildResponse", err)
			span.LogFields(
				opentracinglog.String("event", "buildResponse"),
				opentracinglog.Error(err))
			flow.AbortErr(err)
			ext.Error.Set(span, true)
		} else {
			span.LogFields(opentracinglog.String("event", "ClosedBody"))
		}
		span.LogFields(
			opentracinglog.String("Req", utils.Base64(reqBuf)),
			opentracinglog.String("resp", utils.Base64(respBuf)),
		)

		// cancel nCtx
		if cancel != nil {
			cancel()
		}
	})
}

//func (c *Context) metric() core.Plugin {
//	return core.Function(func(ctx context.Context, flow core.Core) {
//		start := time.Now()
//		flow.Next(ctx)
//		cc := ctx.Value(http_client_context).(*Context)
//		rspCode := http.StatusOK
//		if cc.rsp == nil || cc.core.Err() != nil {
//			rspCode = http.StatusInternalServerError
//		} else {
//			rspCode = cc.rsp.StatusCode
//		}
//		path := cc.Req.raw.URL.Path
//		methodName := strings.Replace(path, "/", ".", -1)
//		methodName = strings.TrimLeft(methodName, ".")
//		metrics.Timer("client."+methodName, start, "clienttag", "httpclient", "httpcode", rspCode)
//	})
//}

func (c *Context) urlParser() core.Plugin {
	return core.Function(func(ctx context.Context, flow core.Core) {
		span := opentracing.SpanFromContext(ctx)
		span.LogFields(opentracinglog.String("event", "URLParsing"))
		cc := ctx.Value(httpClientInternalContext).(*Context)
		err := cc.Req.parseURL()
		if err != nil {
			b := retry.BreakError{
				Err: err,
			}
			flow.AbortErr(b)
			span.LogFields(
				opentracinglog.String("event", "error"),
				opentracinglog.Error(err),
			)
			ext.Error.Set(span, true)
		}
		// inject context to request header
		tracing.ContextToHTTP(ctx, c.client.options.tracer, cc.Req.raw)
	})
}

// use upstream host on http request
func (c *Context) Factory(host string) (core.Plugin, error) {
	return core.Function(func(ctx context.Context, flow core.Core) {
		span := opentracing.SpanFromContext(ctx)
		span.LogFields(opentracinglog.String("event", "DiscoverService"))
		if len(host) == 0 {
			err := fmt.Errorf("host not found from upstream")
			blog := ctx.Value(businessLogger).(log.Logger)
			blog.Log("error", err)
			span.LogFields(
				opentracinglog.String("event", "error"),
				opentracinglog.Error(err),
			)
			ext.Error.Set(span, true)
			flow.AbortErr(err)
		} else {
			cc := ctx.Value(httpClientInternalContext).(*Context)
			cc.host = host
			cc.Req.raw.URL.Host = host
		}
	}), nil
}
