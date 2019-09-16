package memory

import (
	"git.inke.cn/inkelogic/daenerys/config/source"
	"golang.org/x/net/context"
)

type rawChangeSetKey struct{}
type jsonchangeSetkey struct{}
type tomlchangeSetkey struct{}

// WithChangeSet allows a changeset to be set
func WithChangeSet(cs *source.ChangeSet) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, rawChangeSetKey{}, cs)
	}
}

// WithData allows the source data to be set
func WithDataJson(d []byte) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, jsonchangeSetkey{}, &source.ChangeSet{
			Data:   d,
			Format: "json",
		})
	}
}

func WithDataToml(d []byte) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tomlchangeSetkey{}, &source.ChangeSet{
			Data:   d,
			Format: "toml",
		})
	}
}
