package client

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/context"
)

// client interface
type Client interface {
	Call(*Request) (*Response, error)
}
type Func func(*Request) (*Response, error)

func (f Func) Call(req *Request) (*Response, error) {

	return f(req)
}

type client struct {
	client  *http.Client
	options Options
}

func NewClient(opts ...Option) Client {
	c := &client{}
	c.options = newOptions(opts...)
	if c.options.client != nil {
		c.client = c.options.client
	} else {
		c.client = &http.Client{
			Transport: &http.Transport{
				// 表示连接池对每个host的最大链接数量,默认50
				MaxIdleConnsPerHost: c.options.maxIdleConnsPerHost,

				// 表示连接池对所有host的最大链接数量,默认50
				MaxIdleConns: c.options.maxIdleConns,

				Proxy: http.ProxyFromEnvironment,

				// 该函数用于创建http（非https）连接
				DialContext: (&net.Dialer{
					// 表示建立Tcp链接超时时间,默认30s
					Timeout: c.options.dialTimeout,

					// 表示底层为了维持http keepalive状态,每隔多长时间发送Keep-Alive报文
					// 通常要与IdleConnTimeout对应,默认30s
					KeepAlive: c.options.keepAliveTimeout,
					DualStack: true,
				}).DialContext,

				// 连接最大空闲时间,超过这个时间就会被关闭,也即socket在该时间内没有交互则自动关闭连接
				// 该timeout起点是从每次空闲开始计时,若有交互则重置为0,该参数通常设置为分钟级别,默认90s
				IdleConnTimeout: c.options.idleConnTimeout,

				// 限制TLS握手使用的时间
				TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
				ExpectContinueTimeout: defaultExpectContinueTimeout,

				// 表示是否开启http keepalive功能，也即是否重用连接，默认开启(false)
				DisableKeepAlives: c.options.keepAlivesDisable,
			},
			// 此client的请求处理时间,包括建连,重定向,读resp所需的所有时间
			Timeout: c.options.requestTimeout,
		}
	}
	return c
}

func (c *client) Call(r *Request) (*Response, error) {
	if c == nil {
		return nil, fmt.Errorf("http client must be init first")
	}

	ctx := newContext(c, r, &c.options)
	ctx.Ctx = context.WithValue(ctx.Ctx, httpClientInternalContext, ctx)

	// plugins list:
	// recover -> tracing -> logging -> breaker -> common retry ->
	// upstream -> url parse -> sender(last plugin)
	ctx.core.Next(ctx.Ctx)
	return ctx.Resp, ctx.core.Err()
}
