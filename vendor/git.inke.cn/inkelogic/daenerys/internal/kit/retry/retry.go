package retry

import (
	"fmt"
	"strings"

	"git.inke.cn/inkelogic/daenerys/internal/core"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type RetryError struct {
	RawErrors []error
	Final     error
}

func (e RetryError) Error() string {
	var suffix string
	if len(e.RawErrors) > 1 {
		a := make([]string, len(e.RawErrors)-1)
		for i := 0; i < len(e.RawErrors)-1; i++ { // last one is Final
			a[i] = e.RawErrors[i].Error()
		}
		suffix = fmt.Sprintf(" (previously: %s)", strings.Join(a, "; "))
	}
	return fmt.Sprintf("%v%s", e.Final, suffix)
}

type BreakError struct {
	Err error
}

func (b BreakError) Error() string {
	return b.Err.Error()
}

type KeepTrying interface {
	KeepTrying(n int, received error) bool
}

type trymax struct {
	max int
}

func (t trymax) KeepTrying(n int, received error) bool {
	return n < t.max
}

func Retry(max int) core.Plugin {
	return RetryKeepTrying(trymax{max})
}

func RetryKeepTrying(keep KeepTrying) core.Plugin {
	return core.Function(func(ctx context.Context, c core.Core) {
		var final RetryError
		span := opentracing.SpanFromContext(ctx)

		// TODO copy plugin
		for i := 1; ; i++ {
			core := c.Copy()
			core.Next(ctx)
			err := core.Err()
			if err == nil {
				c.AbortErr(nil)
				return
			}
			switch e := err.(type) {
			case BreakError:
				c.AbortErr(e.Err)
				return
			}

			final.RawErrors = append(final.RawErrors, err)
			final.Final = err

			if keep.KeepTrying(i, err) {
				if span != nil {
					span.LogFields(
						log.String("event", "retrying"),
						log.Int("times", i),
					)
				}
				continue
			}
			break
		}

		c.AbortErr(final)
	})
}
