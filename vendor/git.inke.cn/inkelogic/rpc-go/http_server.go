package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.inke.cn/inkelogic/rpc-go/codes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"git.inke.cn/inkelogic/daenerys"

	httpserver "git.inke.cn/inkelogic/daenerys/http/server"

	"sync/atomic"

	"golang.org/x/net/context"

	log "git.inke.cn/BackendPlatform/golang/logging"
)

var (
	nullBody = []byte{}
)

var (
	apiHandlerNilError = errors.New("register api handler is nil")
)

var (
	stopflag          int64
	requestLogBodyOff bool
	globalHTTPServer  atomic.Value
)

func init() {
	stopflag = 0
}

func SetRequestBodyLogOff(off bool) {
	requestLogBodyOff = off
}

func GetRemoteIp(r *http.Request) string {
	for _, h := range []string{"X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := addresses[i]
			if len(ip) > 0 {
				return ip
			}
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

type httpHeaderKey struct{}

var activeHTTPHeaderKey = httpHeaderKey{}

type httprspHeader struct {
	headers map[string]string
}

func (t *httprspHeader) getHeader() map[string]string {
	if t.headers == nil {
		tmp := make(map[string]string)
		return tmp
	}
	return t.headers
}
func (t *httprspHeader) setHeader(key string, value string) {

	if t.headers == nil {
		t.headers = make(map[string]string)
	}
	t.headers[key] = value
}

func SetHttpRspHeader(ctx context.Context, key string, value string) {
	if ctx.Value(activeHTTPHeaderKey) != nil {
		header, _ := ctx.Value(activeHTTPHeaderKey).(*httprspHeader)
		if header != nil {
			header.setHeader(key, value)
		}
	}
}

type JsonApiHandler interface {
	Serve(context.Context, *http.Request) (interface{}, int)
}
type TextApiHandler interface {
	Serve(context.Context, *http.Request) (interface{}, int)
	Text() bool
}

type restfulHandler struct {
	url     string
	c       JsonApiHandler
	bodyLog bool
}

type RestRecord struct {
	RequestURI   string
	RealIP       string
	StatusCode   int
	Method       string
	TimeCost     time.Duration
	RequestID    interface{}
	RemoteAddr   string
	ResponseBody []byte
	RequestBody  []byte
	Logger       *log.Logger
	Extra        interface{}
}

type RestRecorder interface {
	Record(context.Context, *http.Request, *RestRecord)
}

type HTTPServer struct {
	configPath    []string
	configLogBody []string
	httpSrv       httpserver.Server
}

func newHTTPServer(path, logBody []string) *HTTPServer {
	s := &HTTPServer{
		configPath:    path,
		configLogBody: logBody,
	}
	s.httpSrv = daenerys.Default.HTTPServer()
	return s
}

func (s *HTTPServer) handleRequest(ctx *httpserver.Context, rest *restfulHandler) {
	// copy body
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	// run bussiness logic
	rsp, busiCode := rest.c.Serve(ctx.Ctx, ctx.Request)

	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	defer func() {
		ctx.SetBusiCode(int32(busiCode))
	}()

	// add header to response
	header, _ := ctx.Ctx.Value(activeHTTPHeaderKey).(*httprspHeader)
	if header != nil {
		headers := header.getHeader()
		for headerkey, headervalue := range headers {
			ctx.Response.Header().Add(headerkey, headervalue)
		}
	}

	var jsonEnc bool
	switch rest.c.(type) {
	case TextApiHandler:
		switch rsp.(type) {
		case []byte:
			ctx.Response.Write(rsp.([]byte))
		default:
			if rest.c.(TextApiHandler).Text() {
				result := fmt.Sprintf("%v", rsp)
				ctx.Response.WriteString(result)
			} else {
				jsonEnc = true
			}
		}
	default:
		jsonEnc = true
	}
	if jsonEnc {
		encoder := json.NewEncoder(ctx.Response)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(rsp); err != nil {
			http.Error(ctx.Response.Writer(), "Internal Error", http.StatusInternalServerError)
			busiCode = codes.ParseResponseMessage
			ctx.AbortErr(err)
		}
	}
}

func (s *HTTPServer) requestHandler(url string, l JsonApiHandler, logBody bool) httpserver.HandlerFunc {
	return func(ctx *httpserver.Context) {
		if l == nil || len(url) == 0 {
			ctx.AbortErr(apiHandlerNilError)
			return
		}
		r := &restfulHandler{
			url:     url,
			c:       l,
			bodyLog: logBody,
		}
		ctx.Ctx = context.WithValue(ctx.Ctx, activeHTTPHeaderKey, &httprspHeader{headers: make(map[string]string)})
		s.handleRequest(ctx, r)
	}
}

func (s *HTTPServer) restRecordHandler(l JsonApiHandler, logBody bool) httpserver.HandlerFunc {
	return func(ctx *httpserver.Context) {
		timeNow := time.Now()
		ctx.Next()
		if ctx.Err() != nil {
			return
		}

		// copy body
		busiBody, _ := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(busiBody))

		respBody := nullBody
		if logBody {
			respBody = ctx.Response.ByteBody()
		}

		endTime := time.Now()
		remoteAddr := GetRemoteIp(ctx.Request)
		ip, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)
		requestID := getRequestIDString(ctx.Ctx)
		accessInfo := &RestRecord{
			RequestURI:   ctx.Request.URL.RequestURI(),
			RealIP:       remoteAddr,
			StatusCode:   ctx.Response.Status(),
			Method:       ctx.Request.Method,
			TimeCost:     endTime.Sub(timeNow),
			RequestID:    requestID,
			RemoteAddr:   ip,
			ResponseBody: respBody,
			RequestBody:  busiBody,
			Extra:        nil,
			Logger:       log.Log(log.DefaultLoggerName),
		}

		switch l.(type) {
		case RestRecorder:
			l.(RestRecorder).Record(ctx.Ctx, ctx.Request, accessInfo)
		default:
			if requestLogBodyOff {
				accessInfo.RequestBody = nullBody
			} else if len(accessInfo.RequestBody) > 200 {
				accessInfo.RequestBody = accessInfo.RequestBody[:200]
			}
		}
		if len(accessInfo.RequestBody) == 0 {
			accessInfo.RequestBody = nullBody
		}
		body := fmt.Sprintf("%q", accessInfo.RequestBody)
		ctx.LoggingExtra("req_body", body)

		var buff bytes.Buffer
		buff.Write(respBody)
		ctx.LoggingExtra("resp_body", buff.String())

		if accessInfo.Extra != nil {
			extra := fmt.Sprintf("%v", accessInfo.Extra)
			ctx.LoggingExtra("extra", extra)
		}
	}
}

func (s *HTTPServer) RegisterHandler(url string, l JsonApiHandler, logBody bool) error {
	s.httpSrv.ANY(url, s.restRecordHandler(l, logBody), s.requestHandler(url, l, logBody))
	return nil
}

// Register can only be called once
func (s *HTTPServer) Register(ls ...JsonApiHandler) error {
	apiHandlerLen := len(ls)
	configPathLen := len(s.configPath)
	if apiHandlerLen != configPathLen {
		fmt.Printf("[Warning] HTTPServer serve path and handler length not equal.\n")
	}
	min := apiHandlerLen
	if configPathLen > apiHandlerLen {
		min = apiHandlerLen
	}
	for i := 0; i < min; i++ {
		logBody := true
		if i < len(s.configLogBody) {
			if s.configLogBody[i] == "false" {
				logBody = false
			}
		}
		if err := s.RegisterHandler(s.configPath[i], ls[i], logBody); err != nil {
			return err
		}
	}
	return nil
}

func (s *HTTPServer) Serve(port int) error {
	globalHTTPServer.Store(s)
	addr := fmt.Sprintf(":%d", port)
	err := s.httpSrv.Run(addr)
	if err != nil {
		atomic.CompareAndSwapInt64(&stopflag, 0, 1)
	}
	return err
}

func SetHttpStop() {
	if !atomic.CompareAndSwapInt64(&stopflag, 0, 1) {
		return
	}
	srv := globalHTTPServer.Load()
	if srv != nil {
		srv := srv.(*HTTPServer)
		srv.httpSrv.Stop()
	}
}
