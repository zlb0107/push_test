package core

import (
	"golang.org/x/net/context"
)

type ikCore struct {
	plugins []Plugin
	index   int
	err     error
}

func New(ps ...Plugin) Core {
	c := new(ikCore)
	c.index = -1
	c.Use(ps...)
	return c
}

func (c *ikCore) Copy() Core {
	dup := &ikCore{}
	dup.index = c.index
	dup.plugins = append(c.plugins[:0:0], c.plugins...)
	return dup
}

func (c *ikCore) Use(ps ...Plugin) Core {
	c.plugins = append(c.plugins, ps...)
	return c
}

func (c *ikCore) Next(ctx context.Context) {
	c.index++
	for s := len(c.plugins); c.index < s; c.index++ {
		c.plugins[c.index].Do(ctx, c)
	}
}

func (c *ikCore) Abort() {
	c.index = len(c.plugins)
}

func (c *ikCore) AbortErr(err error) {
	c.Abort()
	c.err = err
}

func (c *ikCore) Err() error {
	return c.err
}

func (c *ikCore) IsAborted() bool {
	return c.index >= len(c.plugins)
}
