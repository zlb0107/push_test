package memory

import (
	"bytes"
	"errors"
	"fmt"
	"git.inke.cn/inkelogic/daenerys/config/loader"
	"git.inke.cn/inkelogic/daenerys/config/reader"
	"git.inke.cn/inkelogic/daenerys/config/reader/toml"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"reflect"
	"strings"
	"sync"
	"time"
)

type memory struct {
	exit chan bool
	opts loader.Options

	sync.RWMutex
	snap    *loader.Snapshot
	vals    reader.Values
	sets    []*source.ChangeSet
	sources []source.Source

	idx      int
	watchers map[int]*watcher
	//for dynamic watch
	exchanged chan bool
}

func NewLoader(opts ...loader.Option) loader.Loader {
	options := loader.Options{
		Reader: toml.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	m := &memory{
		exit:      make(chan bool),
		opts:      options,
		watchers:  make(map[int]*watcher),
		sources:   options.Source,
		exchanged: make(chan bool),
	}

	for i, s := range options.Source {
		go m.watch(i, s)
	}
	return m
}

//实现loader.Loader接口
func (m *memory) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}
func (m *memory) Load(sources ...source.Source) error {
	var gerrors []string

	for _, source := range sources {
		set, err := source.Read() //各数据源用自己的格式解析
		if err != nil {
			gerrors = append(gerrors, fmt.Sprintf("error loading source %s: %v", source, err))
			continue
		}
		m.Lock()
		m.sources = append(m.sources, source)
		m.sets = append(m.sets, set)
		idx := len(m.sets) - 1
		m.Unlock()
		go m.watch(idx, source)
	}

	if err := m.reload(); err != nil {
		gerrors = append(gerrors, err.Error())
	}
	if len(gerrors) != 0 {
		return errors.New(strings.Join(gerrors, "\n"))
	}
	return nil
}
func (m *memory) Snapshot() (*loader.Snapshot, error) {
	if m.loaded() {
		m.RLock()
		snap := loader.Copy(m.snap)
		m.RUnlock()
		return snap, nil
	}

	// not loaded, sync
	if err := m.Sync(); err != nil {
		return nil, err
	}

	// make copy
	m.RLock()
	snap := loader.Copy(m.snap)
	m.RUnlock()

	return snap, nil
}
func (m *memory) Sync() error {
	var sets []*source.ChangeSet
	m.Lock()
	var gerr []string
	for _, source := range m.sources {
		ch, err := source.Read()
		if err != nil {
			gerr = append(gerr, err.Error())
			continue
		}
		sets = append(sets, ch)
	}

	// merge sets
	set, err := m.opts.Reader.Merge(sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	vals, err := m.opts.Reader.Values(set)
	if err != nil {
		m.Unlock()
		return err
	}
	m.vals = vals
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   fmt.Sprintf("%d", time.Now().Unix()),
	}
	m.Unlock()
	m.update()
	if len(gerr) > 0 {
		return fmt.Errorf("source loading errors: %s", strings.Join(gerr, "\n"))
	}

	return nil
}
func (m *memory) Watch(keys ...string) (loader.Watcher, error) {
	value, err := m.get(keys...)
	if err != nil {
		return nil, err
	}

	m.Lock()
	w := &watcher{
		exit:    make(chan bool),
		keys:    keys,
		value:   value,
		reader:  m.opts.Reader,
		updates: make(chan reader.Value, 1),
	}

	id := m.idx
	//每个资源对应一个watcher
	m.watchers[id] = w
	m.idx++

	m.Unlock()

	go func() {
		<-w.exit
		m.Lock()
		delete(m.watchers, id)
		m.Unlock()
	}()

	return w, nil
}
func (m *memory) String() string {
	return "memory"
}

func (m *memory) Listen(v interface{}) loader.Refresher {
	switch cc := v.(type) {
	case loader.AutoLoader:
		return m.listen(cc)
	default:
		vv := &loader.Value{}
		vv.Value.Store(v) //原始值
		vv.Tp = reflect.TypeOf(v)
		vv.Format = m.opts.Reader.String()
		return m.listen(vv)
	}
}

///内部方法////
func (m *memory) listen(cc loader.AutoLoader) loader.Refresher {
	go func() {
		for {
			select {
			case <-m.exchanged:
				m.Lock()
				data := m.vals.Bytes()
				m.Unlock()
				cc.Decode(data)
			}
		}
	}()
	return cc
}
func (m *memory) loaded() bool {
	var loaded bool
	m.RLock()
	if m.vals != nil {
		loaded = true
	}
	m.RUnlock()
	return loaded
}
func (m *memory) flush() error {
	// merge sets, 默认用toml格式合并
	set, err := m.opts.Reader.Merge(m.sets...)
	if err != nil {
		return err
	}

	// set values
	m.vals, _ = m.opts.Reader.Values(set)
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   fmt.Sprintf("%d", time.Now().Unix()),
	}
	return nil
}
func (m *memory) reload() error {
	m.Lock()
	if err := m.flush(); err != nil {
		m.Unlock()
		return err
	}
	m.Unlock()
	m.update()
	return nil
}
func (m *memory) update() {
	var watchers []*watcher
	m.RLock()
	for _, w := range m.watchers {
		watchers = append(watchers, w)
	}
	m.RUnlock()

	for _, w := range watchers {
		m.RLock()
		nval := m.vals.Get(w.keys...)
		m.RUnlock()
		select {
		case w.updates <- nval:
		default:
		}
	}
}
func (m *memory) get(path ...string) (reader.Value, error) {
	if !m.loaded() {
		if err := m.Sync(); err != nil {
			return nil, err
		}
	}

	m.Lock()
	defer m.Unlock()

	if m.vals != nil {
		return m.vals.Get(path...), nil
	}

	ch := m.snap.ChangeSet
	v, err := m.opts.Reader.Values(ch)
	if err != nil {
		return nil, err
	}
	m.vals = v
	if m.vals != nil {
		return m.vals.Get(path...), nil
	}
	return nil, errors.New("no values")
}

//todo:
// 同一个配置项在不同的source里都配置了，从其中一个source删除后，该配置项还会存在
// 需要使用者保证同一个配置项只能配置在一个source中
// idx为加载的先后顺序
func (m *memory) watch(idx int, s source.Source) {
	m.Lock()
	//之前load过的数据
	m.sets = append(m.sets, &source.ChangeSet{Source: s.String()})
	m.Unlock()

	//处理数据更新
	watch := func(idx int, s source.Watcher) error {
		for {
			//获取更新
			cs, err := s.Next()
			if err != nil {
				return err
			}

			m.Lock()
			m.sets[idx] = cs
			if err := m.flush(); err != nil {
				return nil
			}

			select {
			case m.exchanged <- true:
			default:
			}
			m.Unlock()
			m.update()
		}
	}

	for {
		//source Watch
		w, err := s.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)
		go func() {
			select {
			case <-done:
			case <-m.exit:
			}
			w.Stop()
		}()

		// block watch
		if err := watch(idx, w); err != nil {
			time.Sleep(time.Second)
		}
		close(done)
		select {
		case <-m.exit:
			return
		default:
		}
	}
}

//实现loader.Watcher接口
type watcher struct {
	exit    chan bool
	keys    []string
	value   reader.Value
	reader  reader.Reader
	updates chan reader.Value
}

func (w *watcher) Next() (*loader.Snapshot, error) {
	for {
		select {
		case <-w.exit:
			return nil, errors.New("watcher stopped")
		case v := <-w.updates:
			if v == nil {
				continue
			}
			if bytes.Equal(w.value.Bytes(), v.Bytes()) {
				continue
			}

			w.value = v

			cs := &source.ChangeSet{
				Data:      v.Bytes(),
				Format:    w.reader.String(),
				Source:    "memory",
				Timestamp: time.Now(),
			}
			cs.Sum()

			return &loader.Snapshot{
				ChangeSet: cs,
				Version:   fmt.Sprintf("%d", time.Now().Unix()),
			}, nil
		}
	}
}
func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
	}
	return nil
}
