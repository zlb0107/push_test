package ratelimit

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/config"
	"git.inke.cn/inkelogic/daenerys/log"
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"path"
	"sync"
	"sync/atomic"
)

type Config struct {
	namespace  string
	clientname string
	limits     map[string]int
}

var globalWatcher watcher
var once sync.Once

func InitWatcher(prefix string, logger log.Logger) {
	once.Do(func() {
		prefix = path.Join(prefix, "ratelimit")
		globalWatcher.logger = log.WithPrefix(
			logger, "time", log.DefaultTimestamp, "component", "ratelimit",
		)
		globalWatcher.limiters.Store(&sync.Map{})

		wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": prefix})
		if err != nil {
			panic(err)
		}
		wp.Handler = globalWatcher.Decode
		go wp.Run(config.ConsulAddr)
		globalWatcher.logger.Log("event", "watch", "path", prefix)
	})
}

type ratelimitConfig struct {
	Name  string `toml:"name"`
	Limit int    `toml:"limits"`
	Open  bool   `toml:"open"`
}

type configList struct {
	List []ratelimitConfig `toml:"ratelimit"`
}

func NewConfig(namespace, clientname string, limits map[string]int) *Config {
	c := &Config{
		namespace:  namespace,
		clientname: clientname,
		limits:     limits,
	}
	return c
}

type watcher struct {
	logger   log.Logger
	limiters atomic.Value
}

func (w watcher) Decode(idx uint64, data interface{}) {
	if data == nil {
		return
	}

	kvs, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	var body []byte
	for _, v := range kvs {
		body = append(body, v.Value...)
	}

	list := configList{}
	if _, err := toml.Decode(string(body), &list); err != nil {
		w.logger.Log("err", err)
		return
	}

	limiters := &sync.Map{}
	for _, cfg := range list.List {
		limt := NewLimiter(Inf, cfg.Limit)
		limt.SetLimit(Inf)
		if cfg.Open {
			limt.SetLimit(Limit(cfg.Limit))
		}
		limiters.Store(cfg.Name, limt)
		w.logger.Log("name", cfg.Name, "open", cfg.Open, "rate", limt.Limit())
	}
	globalWatcher.limiters.Store(limiters)
}

func (c *Config) Limter(name string) *Limiter {
	limiters := globalWatcher.limiters.Load().(*sync.Map)
	nname := fmt.Sprintf("%s.client.%s.%s", c.namespace, c.clientname, name)
	if l, ok := limiters.Load(nname); ok {
		return l.(*Limiter)
	}
	if val, ok := c.limits[name]; ok {
		l := NewLimiter(Limit(val), val)
		l.SetLimit(Limit(val))
		limiters.Store(nname, l)
		return l
	}
	l := NewLimiter(Inf, 100000)
	l.SetLimit(Inf)
	limiters.Store(nname, l)
	return l
}

func (c *Config) LimterWithPeer(name, peer string) *Limiter {
	limiters := globalWatcher.limiters.Load().(*sync.Map)

	key := fmt.Sprintf("%s.server.%s.%s", c.namespace, peer, name)
	if l, ok := limiters.Load(key); ok {
		return l.(*Limiter)
	}

	key = fmt.Sprintf("%s.server-default.%s", c.namespace, name)
	if l, ok := limiters.Load(key); ok {
		return l.(*Limiter)
	}

	if val, ok := c.limits[name]; ok {
		l := NewLimiter(Limit(val), val)
		l.SetLimit(Limit(val))
		limiters.Store(key, l)
		return l
	}

	l, _ := limiters.LoadOrStore(key, NewLimiter(Inf, 1000000))
	return l.(*Limiter)
}
