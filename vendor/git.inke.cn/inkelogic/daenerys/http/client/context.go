package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"golang.org/x/net/context"
)

const httpClientInternalContext = "_http_client_internal_context_"

type Context struct {
	Ctx  context.Context
	Req  *Request
	Resp *Response

	core   core.Core // 用于外部控制流程
	host   string
	client *client
}

// createResponse creates a default http.Response instance.
func createResponse(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Request:    req,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
	}
}

func newContext(httpclient *client, req *Request, opts *Options) *Context {
	nResp := createResponse(req.raw)
	defer nResp.Body.Close()
	resp, _ := BuildResp(req.raw, nResp)
	c := &Context{
		Ctx:    req.ctx,
		Req:    req,
		Resp:   resp,
		core:   core.New(),
		client: httpclient,
	}
	method := req.raw.Method
	path := req.raw.URL.Path
	service := req.ro.serviceName
	caller := req.ro.callerName
	operation := fmt.Sprintf("HTTP Client %s %s", method, path)
	gPlugins := clientInternalThirdPlugin.OnGlobalStage().Stream()
	ps := append([]core.Plugin{
		c.recover(),
		c.tracing(operation),
		metric.SDName(caller),
		c.logging(),
	}, gPlugins...)

	ps = append(ps, c.ratelimit(path))
	ps = append(ps, c.breaker(service, path))

	if r := c.retry(); r != nil {
		ps = append(ps, r)
	}

	if u := c.upstream(); u != nil {
		ps = append(ps, u)
	} else {
		// ignore upstream plugin
		logging.Debugf("upstream plugin nil, service:%s, path:%s", service, path)
	}

	ps = append(ps, c.urlParser())
	rPlugins := clientInternalThirdPlugin.OnRequestStage().Stream()
	dPlugins := clientInternalThirdPlugin.OnWorkDoneStage().Stream()
	ps = append(ps, rPlugins...)
	ps = append(ps, c.sender())
	ps = append(ps, dPlugins...)
	ps = append(ps, namespace.Namespace(opts.namespace))
	c.core.Use(ps...)
	return c
}

func (c *Context) Next() {
	c.core.Next(c.Ctx)
}

func (c *Context) Abort() {
	c.core.Abort()
}

func (c *Context) AbortErr(err error) {
	c.core.AbortErr(err)
}

func (c *Context) Err() error {
	return c.core.Err()
}
