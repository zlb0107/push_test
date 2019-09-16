package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/rpcerror"
	"git.inke.cn/inkelogic/daenerys/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openlog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

func HClient(endpoint string, options ...Option) Client {
	opts := newOptions(options...)
	return newGeneralClient(&hclient{
		opts: opts,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: opts.maxIdleConnsPerHost,
				MaxIdleConns:        opts.maxIdleConns,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   opts.DialTimeout,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				DisableKeepAlives: opts.keepAlivesDisable,
			},
			Timeout: opts.CallOptions.RequestTimeout,
		},
	}, endpoint, opts)
}

type hclient struct {
	opts   Options
	client *http.Client
}

func (r *hclient) Factory(host string) (core.Plugin, error) {
	return core.Function(func(ctx context.Context, c core.Core) {
		var (
			err    error
			span   = opentracing.SpanFromContext(ctx)
			codec  = r.opts.Codec
			rpcctx = ctx.Value(rpcContextKey).(*rpcContext)
		)
		rpcctx.host = host
		elog := log.WithPrefix(
			r.opts.Error,
			"component", "client-http",
			"time", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"trace", log.TraceID(ctx),
		)
		span.SetTag("proto", "rpc/http")

		defer func() {
			if err == nil {
				return
			}
			ext.Error.Set(span, true)
			c.AbortErr(err)
		}()

		span.LogFields(openlog.String("event", "decode"))
		body, err := codec.Encode(rpcctx.Request)
		if err != nil {
			span.LogFields(
				openlog.String("event", "decode error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, fmt.Errorf("encode: %v", err))
			return
		}

		urlhost := "http://" + host + "/" + strings.Replace(rpcctx.Endpoint, ".", "/", -1)
		request, err := http.NewRequest("POST", urlhost, bytes.NewReader(body))
		if err != nil {
			span.LogFields(
				openlog.String("event", "make request error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, err)
			return
		}
		// record request
		reqBuf, _ := httputil.DumpRequest(request, true)

		if err := tracing.ContextToHTTP(ctx, r.opts.Tracer, request); err != nil {
			elog.Log("err", err)
		}
		span.LogFields(openlog.String("event", "transport"))
		response, err := r.client.Do(request)
		if err != nil {
			if e, ok := err.(*url.Error); ok && e.Timeout() {
				rpcctx.Retry()
				err = rpcerror.Error(rpcerror.Timeout, err)
			}
			span.LogFields(
				openlog.String("event", "transport error"),
				openlog.Error(err),
			)
			elog.Log("err", err)
			return
		}

		defer response.Body.Close()

		span.LogFields(openlog.String("event", "read response"))
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			span.LogFields(
				openlog.String("event", "read response error"),
				openlog.Error(err),
			)
			err = rpcerror.Error(rpcerror.Internal, err)
			return
		}

		// record req&resp
		span.LogFields(
			openlog.String("req", utils.Base64(reqBuf)),
			openlog.String("resp", utils.Base64(body)),
		)

		if response.StatusCode != 200 {
			err = rpcerror.HTTP(body)
			span.LogFields(
				openlog.String("event", "response error"),
				openlog.Error(err),
			)
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
