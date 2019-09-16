package consul

import (
	"errors"
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

type watcher struct {
	e           encoder.Encoder
	name        string
	stripPrefix string

	wp   *watch.Plan
	ch   chan *source.ChangeSet
	exit chan bool
}

func newWatcher(key, addr, name, stripPrefix string, e encoder.Encoder) (source.Watcher, error) {
	w := &watcher{
		e:           e,
		name:        name,
		stripPrefix: stripPrefix,
		ch:          make(chan *source.ChangeSet),
		exit:        make(chan bool),
	}

	//新建一个watch plan
	wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": key})
	if err != nil {
		return nil, err
	}

	//有更新时会调用handler处理
	wp.Handler = w.handle

	//启动wach Plan
	// wp.Run is a blocking call and will prevent newWatcher from returning
	go wp.Run(addr)

	w.wp = wp

	return w, nil
}

func (w *watcher) handle(idx uint64, data interface{}) {
	if data == nil {
		return
	}

	kvs, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	//加载数据
	d, err := makeMap(w.e, kvs, w.stripPrefix)
	if err != nil {
		return
	}

	//编码
	b, err := w.e.Encode(d)
	if err != nil {
		return
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    w.e.String(),
		Source:    w.name,
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	w.ch <- cs
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch: //数据有更新返回
		return cs, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		w.wp.Stop()
		close(w.exit)
	}
	return nil
}
