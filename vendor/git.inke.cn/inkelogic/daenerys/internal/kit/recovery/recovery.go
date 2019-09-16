package recovery

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/tpc/inf/metrics"
	"golang.org/x/net/context"
	"os"
	"runtime/debug"
)

func Recovery(logger log.Logger) core.Plugin {
	logger = log.WithPrefix(logger, "plugin", "recovery")
	return core.Function(func(ctx context.Context, c core.Core) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(os.Stderr, "panic: %s\n%s", err, debug.Stack())
				if logger != nil {
					logger.Log("error", err)
				}
				// TODO metric name and stuck
				metrics.Meter("plugin.recovery", 1)
				c.AbortErr(fmt.Errorf("%s", err))
			}
		}()
		c.Next(ctx)
	})
}
