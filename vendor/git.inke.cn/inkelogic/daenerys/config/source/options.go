package source

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"golang.org/x/net/context"
)

type Options struct {
	Encoder encoder.Encoder
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: toml.NewEncoder(), //数据源默认用toml格式解析
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

//设置数据源用何种格式解析,一个source只能有一种encoder
func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}
