package server

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/rpc/codec"
	"github.com/pkg/errors"
	"sync"
)

type Both struct {
	router *router
	binary Server
	http   Server
	exitC  chan error
	startC chan error
}

func BothServer(port int, options ...Option) Server {
	b := &Both{}
	b.router = newRouter()
	b.exitC = make(chan error, 2)
	b.startC = make(chan error, 2)

	ops1 := append(
		options,
		Address(fmt.Sprintf(":%d", port)),
		Codec(codec.NewProtoCodec()),
	)
	b.binary = BinaryServer(ops1...)

	ops2 := append(
		options,
		Address(fmt.Sprintf(":%d", port+1)),
		Codec(codec.NewJSONCodec()),
	)
	b.http = HTTPServer(ops2...)
	b.Use(RatelimitPlugin)
	b.Use(BreakerPlugin)
	return b
}

func (b *Both) NewHandler(handler interface{}, opts ...HandlerOption) Handler {
	return b.router.NewHandler(handler, opts...)
}

func (b *Both) Handle(h Handler) error {
	if err := b.binary.Handle(h); err != nil {
		return err
	}
	if err := b.http.Handle(h); err != nil {
		return err
	}
	return nil
}

func (b *Both) Use(list ...Plugin) Server {
	b.binary.Use(list...)
	b.http.Use(list...)
	return b
}

func (b *Both) Start() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	go func() {
		e := b.binary.Start()
		if e != nil {
			err = errors.Wrap(err, e.Error())
		}
		wg.Done()
	}()
	go func() {
		e := b.http.Start()
		if e != nil {
			err = errors.Wrap(err, e.Error())
		}
		wg.Done()
	}()
	wg.Wait()
	return err
}

func (b *Both) Stop() error {
	var err error
	if e := b.binary.Stop(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	if e := b.http.Stop(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	return err
}
