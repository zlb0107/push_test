package config

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/loader"
	"git.inke.cn/inkelogic/daenerys/config/loader/memory"
	"git.inke.cn/inkelogic/daenerys/config/reader"
	"git.inke.cn/inkelogic/daenerys/config/reader/toml"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"git.inke.cn/inkelogic/daenerys/config/source/file"
	"sync"
	"time"
)

type defaultConfig struct {
	exit chan bool
	opts Options

	sync.RWMutex
	snap *loader.Snapshot
	vals reader.Values

	//固定
	loader loader.Loader //都加载到memory.NewLoader()中
	reader reader.Reader //统一用toml方式读取合并的内存数据
}

func newDefaultConfig(opts ...Option) *defaultConfig {

	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c := &defaultConfig{
		loader: memory.NewLoader(),
		reader: toml.NewReader(),
	}

	c.loader.Load(ops.Source...)
	snap, err := c.loader.Snapshot()
	if err != nil {
		panic(err)
	}
	vals, err := c.reader.Values(snap.ChangeSet)
	if err != nil {
		panic(err)
	}

	c.exit = make(chan bool)
	c.opts = ops
	c.snap = snap
	c.vals = vals

	go c.run()
	return c
}

func (c *defaultConfig) update() error {
	snap, err := c.loader.Snapshot()
	if err != nil {
		return err
	}
	c.Lock()
	defer c.Unlock()

	c.snap = snap
	//toml方式读取内存中的合并数据
	vals, err := c.reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals
	return nil
}

//Config实现:config.Config
func (c *defaultConfig) Load(sources ...source.Source) error {
	//都加载到memory.NewLoader()中
	if err := c.loader.Load(sources...); err != nil {
		return err
	}
	return c.update()

}
func (c *defaultConfig) LoadFile(f string, encoder encoder.Encoder) error {
	return c.Load(file.NewSource(file.WithPath(f), source.WithEncoder(encoder)))
}
func (c *defaultConfig) Sync() error {
	if err := c.loader.Sync(); err != nil {
		return err
	}
	return c.update()
}
func (c *defaultConfig) Listen(v interface{}) loader.Refresher {
	c.Lock()
	defer c.Unlock()
	return c.loader.Listen(v)
}

//Values实现:reader.Values
func (c *defaultConfig) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.vals == nil {
		return []byte{}
	}
	return c.vals.Bytes()
}
func (c *defaultConfig) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()
	if c.vals != nil {
		return c.vals.Get(path...)
	}
	return nil
}
func (c *defaultConfig) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Map()
}
func (c *defaultConfig) Scan(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.vals.Scan(v)
}

//数据变更监视
func (c *defaultConfig) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get changeset
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.Lock()
			c.snap = snap
			c.vals, _ = c.reader.Values(snap.ChangeSet)
			c.Unlock()
		}
	}

	for {
		// memory loader's watcher
		w, err := c.loader.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)
		// the stop watch func
		go func() {
			<-done
			w.Stop()
		}()

		// block watch
		if err := watch(w); err != nil {
			time.Sleep(time.Second)
		}
		close(done)
	}
}
