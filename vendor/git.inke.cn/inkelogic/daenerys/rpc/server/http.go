package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
	"git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"git.inke.cn/inkelogic/daenerys/internal/kit/tracing"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/rpc/internal/rpcerror"
	"git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/tpc/inf/go-tls"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"git.inke.cn/tpc/inf/metrics"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	openlog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type httpServer struct {
	opts     Options
	router   *router
	srv      *http.Server
	plugins  []Plugin
	cfg      *config.Register
	stop     chan struct{}
	shutdown int32
	once     sync.Once
}

func HTTPServer(options ...Option) Server {
	opts := newOptions(options...)
	h := &httpServer{}
	h.router = newRouter()
	h.opts = opts
	h.srv = &http.Server{
		Addr:      opts.Address,
		Handler:   h,
		ConnState: nil, // TODO
	}
	h.stop = make(chan struct{})
	h.shutdown = 0
	return h
}

func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var codec = h.opts.Codec
	var remoteAddr = r.RemoteAddr
	var start = time.Now()

	// record HttpRPC request
	reqBuf, _ := httputil.DumpRequest(r, true)

	elog := log.WithPrefix(
		h.opts.Error,
		"component", "rpcserver-http",
		"time", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"remote", remoteAddr,
	)

	optname := fmt.Sprintf("RPC Server %s %s", r.Method, r.URL.Path)
	ctx := tracing.HTTPToContext(h.opts.Tracer, r, optname, elog)
	tls.SetContext(ctx)
	ns := namespace.GetNamespace(ctx)
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("proto", "rpc/http")
	peer := metric.GetSDName(ctx)
	defer tls.Flush()
	defer span.Finish()

	writeError := func(code int, dest string, name string) {
		ext.Error.Set(span, true)
		http.Error(w, rpcerror.HTTPError{
			C:    code,
			Desc: dest,
		}.Marshal(), http.StatusBadRequest)
		if ns != "" {
			metrics.Timer(name, start, metrics.TagCode, code, "namespace", ns)
		} else {
			metrics.Timer(name, start, metrics.TagCode, code)
		}
	}

	names := strings.Replace(r.URL.Path[1:], "/", ".", -1)
	pos := strings.LastIndex(names, ".")
	if pos == -1 {
		writeError(rpcerror.Internal, fmt.Sprintf("Malfored Method Name"), "HServer")
		span.LogFields(
			openlog.String("event", "decode meta data"),
			openlog.Error(errors.New("Malfored Method Name")),
		)
		elog.Log("err", "Malfored Method Name", "url", r.URL.Path)
		return
	}

	service := names[:pos]
	method := names[pos+1:]
	metricname := fmt.Sprintf("HServer.%s.%s", service, method)

	span.SetOperationName(fmt.Sprintf("RPC Server %s.%s", service, method))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(rpcerror.Internal, "Read Body Error", metricname)
		span.LogFields(
			openlog.String("event", "read body"),
			openlog.Error(err),
		)
		elog.Log("err", err)
		return
	}

	stype, mtype, args, err := h.router.signature(service, method)
	if err != nil {
		writeError(rpcerror.Internal, fmt.Sprintf("Unknown Method %v", method), metricname)
		elog.Log("err", err)
		span.LogFields(openlog.Error(err))
		return
	}

	blog := log.WithPrefix(
		h.opts.Kit.B(),
		"time", log.DefaultTimestamp,
		"end", "server",
		"trace", log.TraceID(ctx),
		"remote", remoteAddr,
		"cost", log.Cost(start),
		"service", fmt.Sprintf("%s.%s", service, method),
		"proto", "http",
		"peer", peer,
	)

	if ns != "" {
		blog = log.WithPrefix(blog, "namespace", ns)
	}

	c := core.New()
	rpcctx := &Context{
		core:       c,
		opts:       h.opts,
		Ctx:        ctx,
		Service:    service,
		Method:     method,
		Peer:       peer,
		Header:     make(map[string]string),
		RemoteAddr: remoteAddr,
		Body:       body,
		Request:    args,
		Code:       int32(rpcerror.Success),
		Namespace:  ns,
	}

	span.LogFields(openlog.String("event", "decode"))
	if len(body) == 0 {
		body = []byte("{}")
	}
	if err := codec.Decode(body, rpcctx.Request); err != nil {
		writeError(rpcerror.Internal, fmt.Sprintf("Parse Error %v", err), metricname)
		blog.Log("err", err, "code", rpcerror.Internal)
		span.LogFields(openlog.Error(err))
		return
	}

	for _, plugin := range h.plugins {
		plugin := plugin

		c.Use(core.Function(func(ctx context.Context, c core.Core) {
			plugin(rpcctx)
		}))
	}

	c.Use(core.Function(func(ctx context.Context, c core.Core) {
		response, err := h.router.call(rpcctx.Ctx, stype, mtype, rpcctx.Request)
		if err != nil {
			c.AbortErr(err) // TODO
			return
		}
		rpcctx.Response = response
	}))

	span.LogFields(openlog.String("event", "call"))
	rpcctx.Next()

	if err := rpcctx.Err(); err != nil {
		writeError(rpcerror.FromUser, fmt.Sprintf("%v", err), metricname)
		blog.Log("err", err, "code", rpcerror.FromUser)
		span.LogFields(openlog.Error(err))
		return
	}

	span.LogFields(openlog.String("event", "encode"))
	body, err = codec.Encode(rpcctx.Response)
	if err != nil {
		blog.Log("err", err, "code", rpcerror.Internal)
		writeError(rpcerror.Internal, fmt.Sprintf("%v", err), metricname)
		span.LogFields(openlog.Error(err))
		return
	}

	span.LogFields(openlog.String("event", "write"))
	if _, err := w.Write(body); err != nil {
		span.LogFields(openlog.Error(err))
		blog.Log("err", err, "code", rpcerror.Internal)
		writeError(rpcerror.Internal, fmt.Sprintf("%v", err), metricname)
		return
	}
	atomic.LoadInt32(&rpcctx.Code)
	blog.Log("err", "<nil>", "code", rpcctx.Code)

	// record HttpRPC req&resp
	span.LogFields(
		openlog.String("req", utils.Base64(reqBuf)),
		openlog.String("resp", utils.Base64(body)),
	)

	if ns != "" {
		metrics.Timer(metricname, start, metrics.TagCode, rpcctx.Code, "namespace", ns)
	} else {
		metrics.Timer(metricname, start, metrics.TagCode, rpcctx.Code)
	}
}

func (h *httpServer) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	return h.router.NewHandler(handler, opts...)
}

func (h *httpServer) Handle(handler Handler) error {
	return h.router.Handle(handler)
}

func (h *httpServer) Use(list ...Plugin) Server {
	h.plugins = append(h.plugins, list...)
	return h
}

func (h *httpServer) Start() error {
	var err error
	h.once.Do(func() {
		ln, e := net.Listen("tcp4", h.opts.Address)
		if e != nil {
			err = e
			return
		}
		addr := strings.Split(ln.Addr().String(), ":")
		port, _ := strconv.Atoi(addr[1])
		cfg, e := utils.Register(
			h.opts.Manager, h.opts.Name, "http", h.opts.Tags, config.LocalIPString(), port)
		if e != nil {
			err = e
			return
		}
		h.cfg = cfg

		err = h.srv.Serve(ln)
		if err != nil {
			if err == http.ErrServerClosed {
				logging.Infof("rpc-http server closed: %v", err)
				err = nil
			}
		}
		logging.Infof("waiting for rpc-http server stop done")
		fmt.Println("waiting for rpc-http server stop done")
		// waiting for stop done
		<-h.stop
	})
	return err
}

func (h *httpServer) Stop() error {
	if !atomic.CompareAndSwapInt32(&h.shutdown, 0, 1) {
		return nil
	}

	defer close(h.stop)

	if m := h.opts.Manager; m != nil {
		m.Deregister()
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if err := h.srv.Shutdown(ctx); err != nil {
		logging.Errorf("gracefully shutdown, err:%v", err)
	}
	cancel()
	return err
}
