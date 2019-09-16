// Package config is an interface for dynamic configuration.
package config

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"git.inke.cn/inkelogic/daenerys/config/loader"
	"git.inke.cn/inkelogic/daenerys/config/reader"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"git.inke.cn/inkelogic/daenerys/config/source/consul"
	"git.inke.cn/inkelogic/daenerys/config/source/file"
	"git.inke.cn/inkelogic/daenerys/config/source/memory"
	"git.inke.cn/inkelogic/daenerys/internal/kit/sd"
)

//配置器
type Config interface {
	reader.Values
	Load(source ...source.Source) error
	LoadFile(f string, encoder encoder.Encoder) error
	Sync() error //同步最新配置
	Listen(interface{}) loader.Refresher
}

//默认设置
var Default = New()

func New(opts ...Option) Config {
	return newDefaultConfig(opts...)
}

func DefaultRemotePath(sdname, path string) string {
	remotePath, _ := sd.RegistryKVPath(sdname, path)
	return remotePath
}

//reader.Values
func Bytes() []byte {
	return Default.Bytes()
}

func Map() map[string]interface{} {
	return Default.Map()
}

func Scan(v interface{}) error {
	return Default.Scan(v)
}

func Get(path ...string) reader.Value {
	return Default.Get(path...)
}

//多数据源加载
func Files(files ...string) error {
	sources := make([]source.Source, len(files))
	for i, f := range files {
		if len(f) == 0 {
			continue
		}
		s := file.NewSource(
			file.WithPath(f),
			source.WithEncoder(toml.NewEncoder()))
		sources[i] = s
	}
	if len(sources) == 0 {
		return nil
	}
	return Default.Load(sources...)
}

func Consul(address string, paths ...string) error {
	sources := make([]source.Source, len(paths))
	for i, p := range paths {
		s := consul.NewSource(
			consul.WithAddress(address),
			consul.WithPrefix(p),
			consul.StripPrefix(true),
			source.WithEncoder(toml.NewEncoder()))
		sources[i] = s
	}
	return Default.Load(sources...)
}

func MemoryJson(mds ...[]byte) error {
	sources := make([]source.Source, len(mds))
	for i, m := range mds {
		s := memory.NewSource(
			memory.WithDataJson(m),
			source.WithEncoder(json.NewEncoder()))
		sources[i] = s
	}
	return Default.Load(sources...)
}

func MemoryToml(mds ...[]byte) error {
	sources := make([]source.Source, len(mds))
	for i, m := range mds {
		s := memory.NewSource(
			memory.WithDataToml(m),
			source.WithEncoder(toml.NewEncoder()))
		sources[i] = s
	}
	return Default.Load(sources...)
}

func Load(source ...source.Source) error {
	return Default.Load(source...)
}

func Sync() error {
	return Default.Sync()
}

func Listen(structPtr interface{}) loader.Refresher {
	return Default.Listen(structPtr)
}
