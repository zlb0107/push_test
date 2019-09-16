// package loader manages loading from multiple sources
package loader

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"git.inke.cn/inkelogic/daenerys/config/reader"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"golang.org/x/net/context"
	"reflect"
	"sync/atomic"
)

type Loader interface {
	Load(...source.Source) error
	Snapshot() (*Snapshot, error)
	Sync() error
	String() string
	Listen(v interface{}) Refresher

	Watch(keys ...string) (Watcher, error)
	Close() error
}

// Watcher, watch sources and returns a merged ChangeSet
type Watcher interface {
	Next() (*Snapshot, error)
	Stop() error
}

// Snapshot is a merged ChangeSet
type Snapshot struct {
	ChangeSet *source.ChangeSet
	Version   string
}

type Options struct {
	Reader  reader.Reader
	Source  []source.Source
	Context context.Context
}

type Option func(o *Options)

// Copy snapshot,not deep copy
func Copy(s *Snapshot) *Snapshot {
	cs := *(s.ChangeSet)
	return &Snapshot{
		ChangeSet: &cs,
		Version:   s.Version,
	}
}

/////动态加载/////

type AutoLoader interface {
	Decode([]byte) error
	Refresher
}

type Refresher interface {
	Load() interface{}
}

type Value struct {
	raw    atomic.Value
	Value  atomic.Value
	Tp     reflect.Type
	Format string
}

func (v *Value) Decode(b []byte) error {
	codec, ok := reader.Encoding[v.Format]
	if !ok {
		codec = toml.NewEncoder()
	}
	ins := reflect.New(v.Tp.Elem()).Interface()
	if err := codec.Decode(b, ins); err != nil {
		return err
	}
	v.raw.Store(b)
	v.Value.Store(ins)
	return nil
}

func (v *Value) Load() interface{} {
	return v.Value.Load()
}
