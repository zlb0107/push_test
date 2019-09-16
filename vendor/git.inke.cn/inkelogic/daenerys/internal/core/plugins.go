package core

import (
	"golang.org/x/net/context"
)

type Plugin interface {
	Do(context.Context, Core)
}

type Function func(context.Context, Core)

func (f Function) Do(ctx context.Context, c Core) { f(ctx, c) }
