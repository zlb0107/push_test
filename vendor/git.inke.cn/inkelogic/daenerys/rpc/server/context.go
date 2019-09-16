package server

import (
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"golang.org/x/net/context"
)

type Plugin func(c *Context)

type Context struct {
	core       core.Core
	opts       Options
	Ctx        context.Context
	Service    string
	Method     string
	RemoteAddr string
	Namespace  string
	Peer       string
	Code       int32

	// rpc request raw header
	Header map[string]string

	// rpc request raw body
	Body     []byte
	Request  interface{}
	Response interface{}
}

func (c *Context) Next() {
	c.core.Next(c.Ctx)
}

func (c *Context) AbortErr(err error) {
	c.core.AbortErr(err)
}

func (c *Context) Err() error {
	return c.core.Err()
}
