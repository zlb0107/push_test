package config

import (
	"git.inke.cn/inkelogic/daenerys/config/source"
	"golang.org/x/net/context"
)

type Option func(o *Options)

type Options struct {
	Source  []source.Source
	Context context.Context
}

// WithSource appends a source to list of sources
func WithSource(s source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s)
	}
}
