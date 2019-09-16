package rpc

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.inke.cn/inkelogic/daenerys"
	"github.com/opentracing/opentracing-go"

	log "git.inke.cn/BackendPlatform/golang/logging"

	"git.inke.cn/inkelogic/rpc-go/codes"
	"git.inke.cn/tpc/inf/go-upstream/circuit"
)

// Codec defines the interface RPC uses to encode and decode messages.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// String returns the name of the Codec implementation. The returned
	// string will be used as part of content type in transmission.
	String() string
}

// rpcError defines the status from an RPC.
type rpcError struct {
	Code codes.Code `json:"code,omitempty"`
	Desc string     `json:"error,omitempty"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.Code, e.Desc)
}

// Code returns the error code for err if it was produced by the rpc system.
// Otherwise, it returns codes.Unknown.
func Code(err error) codes.Code {
	if err == nil {
		return codes.Success
	}
	if e, ok := err.(*rpcError); ok {
		return e.Code
	}
	return codes.UnKnown
}

// ErrorDesc returns the error description of err if it was produced by the rpc system.
// Otherwise, it returns err.Error() or empty string when err is nil.
func ErrorDesc(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*rpcError); ok {
		return e.Desc
	}
	return err.Error()
}

// Errorf returns an error containing an error code and a description;
// Errorf returns nil if c is OK.
func Errorf(c codes.Code, format string, a ...interface{}) error {
	if c == codes.Success {
		return nil
	}
	return &rpcError{
		Code: c,
		Desc: fmt.Sprintf(format, a...),
	}
}

// toRPCErr converts an error into a rpcError.
func toRPCErr(err error) *rpcError {
	switch err.(type) {
	case nil:
		return &rpcError{
			Code: codes.Success,
		}
	case *rpcError:
		return err.(*rpcError)
	case *url.Error:
		return &rpcError{
			Code: codes.Unavailable,
			Desc: err.Error(),
		}
	default:
		switch err {
		case io.EOF:
			return &rpcError{
				Code: codes.ChannelBroken,
				Desc: err.Error(),
			}
		case io.ErrClosedPipe, io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite, io.ErrUnexpectedEOF:
			return &rpcError{
				Code: codes.FailedPrecondition,
				Desc: err.Error(),
			}
		case context.DeadlineExceeded:
			return &rpcError{
				Code: codes.DeadlineExceeded,
				Desc: err.Error(),
			}
		case context.Canceled:
			return &rpcError{
				Code: codes.Canceled,
				Desc: err.Error(),
			}
		case ErrClientConnClosing:
			return &rpcError{
				Code: codes.FailedPrecondition,
				Desc: err.Error(),
			}
		case ErrClientLB:
			return &rpcError{
				Code: codes.ConfigLb,
				Desc: err.Error(),
			}
		case circuit.ErrMaxConcurrent, circuit.ErrRateLimit, circuit.ErrSystemLoad:
			return &rpcError{
				Code: codes.BreakerOpen,
				Desc: err.Error(),
			}
		case circuit.ErrAverageRT, circuit.ErrConsecutive, circuit.ErrPercent:
			return &rpcError{
				Code: codes.BreakerOpen,
				Desc: err.Error(),
			}
		case circuit.ErrOpen:
			return &rpcError{
				Code: codes.BreakerOpen,
				Desc: err.Error(),
			}
		}
	}
	return Errorf(codes.UnKnown, "%v", err).(*rpcError)
}

func doExitWork(server *Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		curTime := time.Now().Format("2006-01-02 15:04:05.999")
		SetHttpStop()
		if server != nil {
			server.Stop()
		}
		daenerys.Shutdown()
		log.Infof("%s Shutdown....", curTime)
		fmt.Printf("%s Shutdown...\n", curTime)
		SyncLog()
		os.Exit(0)
	}()
}

func doExitWorkHttpserver(server *HTTPServer) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		curTime := time.Now().Format("2006-01-02 15:04:05.999")
		SetHttpStop()
		if server != nil {
			server.httpSrv.Stop()
		}
		daenerys.Shutdown()
		log.Infof("%s Shutdown....", curTime)
		fmt.Printf("%s Shutdown...\n", curTime)
		SyncLog()
		os.Exit(0)
	}()
}

func NewServerWithConfig(c Config) *Server {
	log.GenLog("start service port:", c.Port(), ",http port:", c.Port()+1)
	server := NewServer()
	doExitWork(server)
	return server
}

func NewHTTPServerWithConfig(c Config) *HTTPServer {
	log.GenLog("start service ,http port:", c.Port())
	SetRequestBodyLogOff(c.RequestBodyLogOff())
	server := newHTTPServer(strings.Split(c.HTTPServeLocation(), ","),
		strings.Split(c.HTTPServeLogBody(), ","))
	doExitWorkHttpserver(server)
	return server
}

// GetRequestID from context
func GetRequestID(ctx context.Context) interface{} {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	return strings.SplitN(fmt.Sprintf("%s", span.Context()), ":", 2)[0]
}

func getRequestIDString(ctx context.Context) string {
	var reqID string
	if traceID := GetRequestID(ctx); traceID != nil {
		reqID = fmt.Sprintf("%s", traceID)
	}
	return reqID
}
