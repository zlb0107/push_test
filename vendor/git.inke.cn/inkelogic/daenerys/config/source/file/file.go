// Package file is a file source. Expected format is json
package file

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"io/ioutil"
	"os"
	"strings"
)

type file struct {
	path string
	opts source.Options
}

//实现Source接口
func (f *file) Read() (*source.ChangeSet, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    format(f.path, f.opts.Encoder),
		Source:    f.String(),
		Timestamp: info.ModTime(),//touch,vim都会更新
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}
func (f *file) String() string {
	return "file"
}
func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	//file watcher
	return newWatcher(f)
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	var path string
	f, ok := options.Context.Value(filePathKey{}).(string)
	if ok {
		path = f
	}
	return &file{opts: options, path: path}
}

func format(p string, e encoder.Encoder) string {
	parts := strings.Split(p, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return e.String()
}