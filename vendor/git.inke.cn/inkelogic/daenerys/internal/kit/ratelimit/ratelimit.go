package ratelimit

import (
	"errors"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	"golang.org/x/net/context"
)

var ErrLimited = errors.New("rate limit exceeded")

// Allower dictates whether or not a request is acceptable to run.
// The Limiter from "git.inke.cn/inkelogic/daenerys/ratelimit" already implements this interface,
// one is able to use that in NewLimiter without any modifications.
type Allower interface {
	Allow() bool
}

func Limiter(name string, config *ratelimit.Config) core.Plugin {
	return core.Function(func(ctx context.Context, c core.Core) {
		if config != nil && !config.Limter(name).Allow() {
			// TODO
			c.AbortErr(ErrLimited)
			return
		}
		c.Next(ctx)
	})
}
