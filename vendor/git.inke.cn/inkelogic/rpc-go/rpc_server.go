package rpc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"git.inke.cn/inkelogic/daenerys"
	rpcserver "git.inke.cn/inkelogic/daenerys/rpc/server"
	"git.inke.cn/inkelogic/rpc-go/desc"
)

var ParseRequestDataError = errors.New("parse request data error")

type dests struct {
	desc *desc.ServiceDesc
	ss   interface{}
}

type Server struct {
	s    rpcserver.Server
	desc []dests
}

// NewServer creates a RPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(...interface{}) *Server {
	srv := &Server{
		desc: make([]dests, 0),
	}
	return srv
}

// RegisterService register a service and its implementation to the RPC
// server. Called from the IDL generated code. This must be called before
// invoking Serve.
func (s *Server) RegisterService(sd *desc.ServiceDesc, ss interface{}) {
	s.desc = append(s.desc, dests{sd, ss})
}

var (
	// ErrServerStopped indicates that the operation is now illegal because of
	// the server being stopped.
	ErrServerStopped = errors.New("rpc: the server has been stopped")
)

func (s *Server) Unregister() {}

func GetDoneWaitGroup(ctx context.Context) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond * 5)
		wg.Done()
	}()
	return wg
}

var activePayLoadKey = "p1"

func GetPayload(ctx context.Context) interface{} {
	return ctx.Value(activePayLoadKey)
}

var activeRemoteAddressKey = "f1"

func GetRemoteAddr(ctx context.Context) interface{} {
	return ctx.Value(activeRemoteAddressKey)
}

var activeRespCodeKey = "c1"

func SetResponseCode(ctx context.Context, code int32) {
	if ctx.Value(activeRespCodeKey) != nil {
		atomic.StoreInt32(ctx.Value(activeRespCodeKey).(*int32), code)
	}
}

// Serve start server client
func (s *Server) Serve(port int) error {
	s.s = daenerys.RPCServer()
	for _, d := range s.desc {
		handler := s.s.NewHandler(d.ss, rpcserver.HandlerName(d.desc.ServiceName))
		if err := s.s.Handle(handler); err != nil {
			return err
		}
	}
	s.s.Use(func(c *rpcserver.Context) {
		c.Ctx = context.WithValue(c.Ctx, activePayLoadKey, c.Body)
		c.Ctx = context.WithValue(c.Ctx, activeRemoteAddressKey, c.RemoteAddr)
		c.Ctx = context.WithValue(c.Ctx, activeRespCodeKey, &c.Code)
		c.Next()
	})
	return s.s.Start()
}

// Stop close the server
func (s *Server) Stop() { s.s.Stop() }
