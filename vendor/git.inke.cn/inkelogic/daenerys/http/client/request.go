package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/context"
)

type RequestOption struct {
	retryTimes  int
	reqTimeout  int // ms
	slowTime    int // ms
	serviceName string
	callerName  string
}

func (ro *RequestOption) RetryTimes(cnt int) *RequestOption {
	ro.retryTimes = cnt
	return ro
}

func (ro *RequestOption) RequestTimeoutMS(timeout int) *RequestOption {
	ro.reqTimeout = timeout
	return ro
}

func (ro *RequestOption) SlowTimeoutMS(timeout int) *RequestOption {
	ro.slowTime = timeout
	return ro
}

func (ro *RequestOption) ServiceName(name string) *RequestOption {
	ro.serviceName = name
	return ro
}

func (ro *RequestOption) CallerName(name string) *RequestOption {
	ro.callerName = name
	return ro
}

type Request struct {
	raw        *http.Request
	ctx        context.Context
	pathParams map[string]string
	queryParam url.Values
	ro         *RequestOption
}

func createRequest() *http.Request {
	req := &http.Request{
		Method:     "GET",
		URL:        &url.URL{Scheme: "http"},
		Host:       "",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Proto:      "HTTP/1.1",
		Header:     make(http.Header),
	}
	return req
}

func NewRequest() *Request {
	r := &Request{
		raw:        createRequest(),
		ctx:        context.Background(),
		queryParam: url.Values{},
		pathParams: map[string]string{},
		ro:         &RequestOption{},
	}
	return r
}

func (r *Request) WithOption(ro *RequestOption) *Request {
	r.ro = ro
	return r
}

func (r *Request) WithServiceName(serviceName string) *Request {
	if r.ro == nil {
		r.ro = &RequestOption{}
	}
	r.ro.serviceName = serviceName
	return r
}

func (r *Request) WithRequest(req *http.Request) *Request {
	r.raw = req
	return r
}

func (r *Request) WithCtxInfo(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) WithMethod(name string) *Request {
	r.raw.Method = name
	return r
}

func normalize(uri string) string {
	match, _ := regexp.MatchString("^http[s]?://", uri)
	if match {
		return uri
	}
	return "http://" + uri
}

func (r *Request) WithURL(uri string) *Request {
	u, err := url.Parse(normalize(uri))
	if err != nil {
		panic("invalid uri")
	}
	r.raw.URL = u
	return r
}

func (r *Request) WithScheme(scheme string) *Request {
	if scheme != "http" && scheme != "https" {
		return r
	}
	r.raw.URL.Scheme = scheme
	return r
}

func (r *Request) WithPath(path string) *Request {
	if path == "/" {
		r.raw.URL.Path = ""
		return r
	}
	r.raw.URL.Path = path
	return r
}

func (r *Request) WithBody(body io.Reader) *Request {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.raw.Body = rc
	// 设置length长度和GetBody函数
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			r.raw.ContentLength = int64(v.Len())
			buf := v.Bytes()
			r.raw.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			r.raw.ContentLength = int64(v.Len())
			snapshot := *v
			r.raw.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			r.raw.ContentLength = int64(v.Len())
			snapshot := *v
			r.raw.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
		}
		if r.raw.GetBody != nil && r.raw.ContentLength == 0 {
			r.raw.Body = http.NoBody
			r.raw.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}

	return r
}

func (r *Request) WithCookie(ck *http.Cookie) *Request {
	r.raw.AddCookie(ck)
	return r
}

func (r *Request) WithMultiCookie(cks []*http.Cookie) *Request {
	for _, v := range cks {
		vv := v
		r.raw.AddCookie(vv)
	}
	return r
}

func (r *Request) WithMultiHeader(headers map[string]string) *Request {
	for k, v := range headers {
		r.raw.Header.Set(k, v)
	}
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.raw.Header.Add(key, value)
	return r
}

func (r *Request) DelHeader(key string) *Request {
	r.raw.Header.Del(key)
	return r
}

// url /v1/user?a=b&c=d
func (r *Request) WithQueryParam(param, value string) *Request {
	r.queryParam.Set(param, value)
	return r
}

func (r *Request) WithMultiQueryParam(params map[string]string) *Request {
	for p, v := range params {
		r.queryParam.Set(p, v)
	}
	return r
}

// `application/x-www-form-urlencoded`
func (r *Request) WithFormData(data map[string]string) *Request {
	for k, v := range data {
		r.raw.Form.Set(k, v)
	}
	return r
}

// url params: /v1/users/:userId/:subAccountId/details
func (r *Request) WithPathParams(params map[string]string) *Request {
	for p, v := range params {
		r.pathParams[p] = v
	}
	return r
}

func (r *Request) RawRequest() *http.Request {
	return r.raw
}

func (r *Request) parseURL() error {
	if len(r.raw.URL.Host) == 0 && len(r.raw.Host) == 0 {
		return fmt.Errorf("request host empty: %s", r.raw.URL.String())
	}
	for p, v := range r.pathParams {
		r.raw.URL.Path = strings.Replace(r.raw.URL.Path, ":"+p, url.PathEscape(v), -1)
	}
	// Adding Query Param
	query := r.raw.URL.Query()
	for k, v := range r.queryParam {
		query.Del(k)
		for _, iv := range v {
			query.Add(k, iv)
		}
	}
	r.raw.URL.RawQuery = query.Encode()
	return nil
}
