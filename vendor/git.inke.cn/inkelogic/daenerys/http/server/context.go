package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/utils"
	"golang.org/x/net/context"
)

const httpServerInternalContext = "_http_server_internal_context_"

type Context struct {
	Request      *http.Request
	Response     Responser
	Params       Params
	Path         string // raw match path
	Peer         string
	Ctx          context.Context // for trace or others store
	opts         *Options
	Namespace    string
	core         core.Core // a control flow
	w            responseWriter
	busiCode     int32
	loggingExtra map[string]interface{}
	bodyBuff     *bytes.Buffer
}

func newContext(w http.ResponseWriter, r *http.Request, opts *Options) *Context {
	ctx := &Context{
		core:         core.New(),
		Ctx:          context.Background(),
		opts:         opts,
		busiCode:     0,
		loggingExtra: map[string]interface{}{},
		bodyBuff:     bytes.NewBuffer([]byte{}),
	}
	var bodyByte []byte
	if r != nil && r.Body != nil {
		bodyByte, _ = ioutil.ReadAll(r.Body)
		ctx.bodyBuff = bytes.NewBuffer(bodyByte)
		r.Body = ioutil.NopCloser(ctx.bodyBuff)
	}
	ctx.Request = r
	ctx.w.reset(w)
	ctx.Response = &ctx.w
	return ctx
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

func (c *Context) SetBusiCode(code int32) {
	atomic.StoreInt32(&c.busiCode, code)
}

func (c *Context) BusiCode() int32 {
	return atomic.LoadInt32(&c.busiCode)
}

func (c *Context) LoggingExtra(vals ...interface{}) {
	if len(vals)%2 != 0 {
		vals = append(vals, log.ErrMissingValue)
	}
	size := len(vals)
	for i := 0; i < size; i += 2 {
		key := fmt.Sprintf("%s", vals[i])
		c.loggingExtra[key] = vals[i+1]
	}
}

func (c *Context) Bind(r *http.Request, model interface{}) error {
	return utils.Bind(r, model)
}

// write response, error include business code and error msg
func (c *Context) JSON(data interface{}, err error) {
	c.Response.WriteHeader(c.Response.Status())
	w := utils.NewWrapResp(data, err)
	c.SetBusiCode(int32(w.Code))
	c.Response.WriteJSON(w)
}

// wrap on JSON
func (c *Context) JSONAbort(data interface{}, err error) {
	c.JSON(data, err)
	c.Abort()
}

// write response, data include error info
func (c *Context) Raw(data interface{}, code int32) {
	c.SetBusiCode(code)
	c.Response.WriteJSON(data)
}
