package server

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/internal/kit/ratelimit"
)

func BreakerPlugin(c *Context) {
	if c.opts.Breaker == nil {
		c.Next()
		return
	}
	endpoint := fmt.Sprintf("%s.%s", c.Service, c.Method)
	brk := c.opts.Breaker.BreakerServer(endpoint)
	if brk != nil {
		err := brk.Call(func() error {
			c.Next()
			return c.Err()
		}, 0)
		if err != nil {
			c.AbortErr(err)
		}
	} else {
		c.Next()
	}
}

func RatelimitPlugin(c *Context) {
	if c.opts.Ratelimit == nil {
		c.Next()
		return
	}
	endpoint := fmt.Sprintf("%s.%s", c.Service, c.Method)
	limt := c.opts.Ratelimit.LimterWithPeer(endpoint, c.Peer)
	if limt != nil && !limt.Allow() {
		c.AbortErr(ratelimit.ErrLimited)
		return
	}
	c.Next()
}
