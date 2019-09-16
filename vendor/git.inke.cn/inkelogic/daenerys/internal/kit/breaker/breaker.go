package breaker

import (
	"golang.org/x/net/context"

	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/internal/kit/metric"
)

type ikBreaker struct {
	name string
	brk  *breaker.Config
}

type ikBreakerServer struct {
	name string
	brk  *breaker.Config
}

// TODO failback
func Breaker(name string, brk *breaker.Config) core.Plugin {
	return ikBreaker{name, brk}
}

func (i ikBreaker) Do(ctx context.Context, c core.Core) {
	if i.brk == nil {
		c.Next(ctx)
		return
	}
	err := i.brk.Breaker(i.name).Call(func() error {
		c.Next(ctx)
		return c.Err()
	}, 0)
	if err != nil {
		c.AbortErr(metric.Error(-2, err))
	}
}

func BreakerServer(name string, brk *breaker.Config) core.Plugin {
	return ikBreakerServer{name, brk}
}

func (i ikBreakerServer) Do(ctx context.Context, c core.Core) {
	if i.brk == nil {
		c.Next(ctx)
		return
	}

	err := i.brk.BreakerServer(i.name).Call(func() error {
		c.Next(ctx)
		return c.Err()
	}, 0)
	if err != nil {
		c.AbortErr(metric.Error(-2, err))
	}
}
