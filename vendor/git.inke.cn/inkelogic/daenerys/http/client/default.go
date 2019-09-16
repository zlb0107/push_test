package client

import (
	"bytes"
	"net"
	"net/http"
)

var DefaultClient Client

var defaultHttpClient *http.Client

func init() {
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
			MaxIdleConns:        defaultMaxIdleConns,
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout,
				KeepAlive: defaultKeepAliveTimeout,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       defaultIdleConnTimeout,
			TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
			ExpectContinueTimeout: defaultExpectContinueTimeout,
			DisableKeepAlives:     false,
		},
		Timeout: defaultRequestTimeout,
	}

	DefaultClient = NewClient(WithClient(defaultHttpClient))
}

func HTTPGet(url string) (*Response, error) {
	r := NewRequest().WithMethod(MethodGet).WithURL(url)
	return DefaultClient.Call(r)
}

func HTTPPost(url string, body []byte) (*Response, error) {
	bf := bytes.NewBuffer(body)
	r := NewRequest().WithMethod(MethodPost).WithURL(url).WithBody(bf)
	return DefaultClient.Call(r)
}

func HTTPPut(url string) (*Response, error) {
	r := NewRequest().WithMethod(MethodPut).WithURL(url)
	return DefaultClient.Call(r)
}

func HTTPDelete(url string) (*Response, error) {
	r := NewRequest().WithMethod(MethodDelete).WithURL(url)
	return DefaultClient.Call(r)
}
