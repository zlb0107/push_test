package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.inke.cn/BackendPlatform/golang/logging"

	"golang.org/x/net/context"

	"net/http/httputil"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	dutils "git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/tpc/inf/go-tls"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

const (
	BreakerOpen = 105
)

type Server interface {
	Router
	Run(addr ...string) error
	Stop() error
	Use(p ...HandlerFunc)
}

// ******** server instance ********//
type server struct {
	RouterMgr
	options  Options
	pluginMu sync.Mutex
	plugins  []core.Plugin
	// flow           core.Core
	trees          methodTrees
	srv            *http.Server
	registryConfig *config.Register
	running        int32
	stop           chan struct{}
	once           sync.Once
}

func NewServer(options ...Option) Server {
	s := &server{
		RouterMgr: RouterMgr{
			plugins:  nil,
			basePath: "/",
		},
		trees:    make(methodTrees, 0, 9),
		pluginMu: sync.Mutex{},
		plugins:  make([]core.Plugin, 0),
		// flow:  core.New(),
		stop: make(chan struct{}),
	}
	s.options = newOptions(options...)
	s.srv = &http.Server{
		Handler:      s,
		ReadTimeout:  s.options.readTimeout,
		WriteTimeout: s.options.writeTimeout,
		IdleTimeout:  s.options.idleTimeout,
	}
	// recover -> logging -> tracing
	s.plugins = []core.Plugin{s.recover(), s.logging(), s.tracing()}
	// s.flow.Use([]core.Plugin{s.recover(), s.logging(), s.tracing()}...)
	s.RouterMgr.server = s
	atomic.StoreInt32(&s.running, 0)
	s.Use(ratelimitPlugin)
	return s
}

func (s *server) Run(addr ...string) error {
	var err error
	s.once.Do(func() {
		port := 0
		if len(addr) > 0 {
			s.srv.Addr = addr[0]
			tmp := strings.Split(s.srv.Addr, ":")
			if len(tmp) == 2 {
				port, _ = strconv.Atoi(tmp[1])
			} else {
				err = fmt.Errorf("invalid addr(s): %v", addr)
			}
		} else if s.options.port > 0 {
			port = s.options.port
			s.srv.Addr = fmt.Sprintf(":%d", port)
		} else {
			s.srv.Addr = ":80"
		}
		ln, e := net.Listen("tcp", s.srv.Addr)
		if e != nil {
			err = e
			return
		}
		logging.Infof("start http server on %s", s.srv.Addr)
		fmt.Printf("start http server on %s\n", s.srv.Addr)
		var cfg *config.Register
		cfg, err = dutils.Register(s.options.manager, s.options.serviceName, "http", s.options.tags, config.LocalIPString(), port)
		if err != nil {
			return
		}
		s.registryConfig = cfg

		if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
			err = fmt.Errorf("server had been running")
			return
		}
		if len(s.options.certFile) == 0 || len(s.options.keyFile) == 0 {
			err = s.srv.Serve(ln)
		} else {
			err = s.srv.ServeTLS(ln, s.options.certFile, s.options.keyFile)
		}
		if err != nil {
			if err == http.ErrServerClosed {
				logging.Infof("http server closed: %v", err)
				err = nil
			}
		}
		logging.Infof("waiting for http server stop done")
		fmt.Println("waiting for http server stop done")
		// waiting for stop done
		<-s.stop
	})
	return err
}

func (s *server) Stop() error {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		return nil
	}

	defer close(s.stop)

	if s.options.manager != nil {
		s.options.manager.Deregister()
	}

	// gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := s.srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		logging.Errorf("gracefully shutdown, err:%v", err)
	}
	cancel()
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBuf, _ := httputil.DumpRequest(r, true)
	ctx := newContext(w, r, &s.options)
	operation := fmt.Sprintf("HTTP Server %s %s", r.Method, r.URL.Path)
	ctx.Ctx = tracing.HTTPToContext(s.options.tracer, r, operation, log.Stdout())
	ctx.Namespace = namespace.GetNamespace(ctx.Ctx)
	ctx.Peer = metric.GetSDName(ctx.Ctx)
	span := opentracing.SpanFromContext(ctx.Ctx)
	tls.SetContext(ctx.Ctx)
	defer tls.Flush()
	defer span.Finish()

	code := http.StatusOK
	flow := s.flow(ctx)
	if flow == nil {
		if s.methodNotAllowed(ctx) {
			span.LogFields(opentracinglog.String("event", "method not allowed"))
			ctx.Response.WriteHeader(http.StatusMethodNotAllowed)
			ctx.Response.Write(default405Body)
			ctx.Response.DoFlush()
			return
		}
		span.LogFields(opentracinglog.String("event", "handlers not found"))
		ctx.Response.WriteHeader(http.StatusNotFound)
		ctx.Response.Write(default404Body)
		ctx.Response.DoFlush()
		return
	}

	// same request use with a copy flow
	ctx.core = flow.Copy()
	ctx.Ctx = context.WithValue(ctx.Ctx, httpServerInternalContext, ctx)
	ctx.core.Next(ctx.Ctx)
	span.LogFields(opentracinglog.String("event", "gotResponse"))
	if err := ctx.core.Err(); err != nil {
		ext.Error.Set(span, true)
		code = http.StatusInternalServerError
	}
	if ctx.Response != nil {
		code = ctx.Response.Status()
	}
	span.LogFields(
		opentracinglog.String("req", dutils.Base64(reqBuf)),
		opentracinglog.String("resp", dutils.Base64(ctx.Response.ByteBody())),
	)
	ctx.Response.WriteHeader(code)
	ctx.Response.DoFlush()
	span.LogFields(opentracinglog.String("event", "wroteResponse"))
}

func (s *server) Use(p ...HandlerFunc) {
	s.pluginMu.Lock()
	defer s.pluginMu.Unlock()
	ps := make([]core.Plugin, len(p))
	for i := range p {
		ps[i] = p[i]
	}
	s.plugins = append(s.plugins, ps...)
	// s.flow.Use(p)
}

func (s *server) addRoute(method, path string, handlers []core.Plugin) {
	if path[0] != '/' || len(method) == 0 || len(handlers) == 0 {
		return
	}
	root := s.trees.get(method)
	if root == nil {
		root = new(node)
		s.trees = append(s.trees, methodTree{method: method, root: root})
	}

	ps := s.makeChain(path, handlers)
	flow := core.New(ps...)
	root.addRoute(path, flow)
}

func (s *server) makeChain(path string, handlers []core.Plugin) []core.Plugin {
	// plugins list:
	// server: recover -> log -> tracing -> (maybe others on server)
	// request plugin on register: breaker -> metric -> handlers
	var ps []core.Plugin
	gPlugins := serverInternalThirdPlugin.OnGlobalStage().Stream()
	ps = append(ps, gPlugins...)

	// plugins on every path: metric -> breaker ->(outside frame plugin) -> handlers -> (outside frame plugin)
	ps = append(ps, s.breaker(path), s.metric())

	rPlugins := serverInternalThirdPlugin.OnRequestStage().Stream()
	dPlugins := serverInternalThirdPlugin.OnWorkDoneStage().Stream()
	ps = append(ps, rPlugins...)
	ps = append(ps, handlers...)
	ps = append(ps, dPlugins...)

	s.pluginMu.Lock()
	defer s.pluginMu.Unlock()
	return append(s.plugins, ps...)
}

func getRemoteIP(r *http.Request) string {
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
