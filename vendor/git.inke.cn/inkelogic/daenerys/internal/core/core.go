package core

import (
	"golang.org/x/net/context"
)

type Core interface {
	Use(...Plugin) Core
	Next(context.Context)
	AbortErr(error)
	Abort()
	IsAborted() bool
	Err() error
	Copy() Core
}
