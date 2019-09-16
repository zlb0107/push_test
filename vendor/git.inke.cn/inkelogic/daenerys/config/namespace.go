package config

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"git.inke.cn/inkelogic/daenerys/config/source/consul"
	"git.inke.cn/inkelogic/daenerys/config/source/file"
	"path/filepath"
)

var ConsulAddr string = "127.0.0.1:8500"

type Namespace struct {
	namespace string
}

func NewNamespace(namespace string) *Namespace {
	return &Namespace{namespace}
}

func NewNamespaceD() *Namespace {
	return &Namespace{"default"}
}

func (m *Namespace) Get(resouce string) Config {
	return m.GetD(resouce, "")
}

func (m *Namespace) With(resouce string) *Namespace {
	return NewNamespace(filepath.Join(m.namespace, resouce))
}

func (m *Namespace) GetD(resouce, filename string) Config {
	return New(
		WithSource(
			file.NewSource(source.WithEncoder(toml.NewEncoder()), file.WithPath(filename)),
		),
		WithSource(
			consul.NewSource(
				consul.WithAddress(ConsulAddr),
				consul.WithPrefix(filepath.Join(m.namespace, resouce)),
				consul.StripPrefix(true),
				source.WithEncoder(toml.NewEncoder()),
			),
		),
	)
}
