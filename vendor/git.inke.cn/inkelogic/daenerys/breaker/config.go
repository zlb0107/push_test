package breaker

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
}

type watcher struct {
	logger   log.Logger
	breakers atomic.Value
}

var globalWatcher watcher
var once sync.Once

type consulConfig struct {
	Breakers []struct {
		Name                      string `toml:"name"`
		ErrorPercentThreshold     int    `toml:"error_percent_threshold"`
		ConsecutiveErrorThreshold int    `toml:"connsecutive_error_threshold"`
		MinSamples                int    `toml:"minsamples"`
		Break                     bool   `toml:"break"`
	} `toml:"breakers"`
}

func InitWatcher(prefix string, logger log.Logger) {
	once.Do(func() {
		prefix = path.Join(prefix, "breaker")
		globalWatcher.logger = log.WithPrefix(
			logger, "time", log.DefaultTimestamp, "component", "breaker",
		)
		globalWatcher.breakers.Store(&sync.Map{})

		wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": prefix})
		if err != nil {
			panic(err)
		}
		wp.Handler = globalWatcher.Decode
		go wp.Run(config.ConsulAddr)
		globalWatcher.logger.Log("event", "watch", "path", prefix)
	})
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

	list := consulConfig{}
	if _, err := toml.Decode(string(body), &list); err != nil {
		w.logger.Log("err", err)
		return
	}
	breakers := &sync.Map{}
	for _, brkcfg := range list.Breakers {
		c := BreakerConfig{
			ErrorPercentThreshold:     brkcfg.ErrorPercentThreshold,
			ConsecutiveErrorThreshold: brkcfg.ConsecutiveErrorThreshold,
			MinSamples:                brkcfg.MinSamples,
			Break:                     brkcfg.Break,
		}
		brk := NewBreakerWithOptions(&Options{
			Name: brkcfg.Name,
		})
		Configure(brkcfg.Name, &c)
		if c.Break {
			brk.Break()
		}
		w.logger.Log(
			"name", brkcfg.Name, "break", c.Break,
			"error_percent", c.ErrorPercentThreshold,
			"connsecutive_error", c.ConsecutiveErrorThreshold,
			"minsamples", c.MinSamples,
		)
		breakers.Store(brkcfg.Name, brk)
	}
	globalWatcher.breakers.Store(breakers)
}

func NewConfig(namespace, clientname string, configs map[string]*BreakerConfig) *Config {
	c := &Config{
		namespace:  namespace,
		clientname: clientname,
	}
	for name, cs := range configs {
		Configure(name, cs)
	}
	return c
}

func (c *Config) Breaker(name string) *Breaker {
	breaker := globalWatcher.breakers.Load().(*sync.Map)
	nname := fmt.Sprintf("%s.client.%s.%s", c.namespace, c.clientname, name)
	if b, ok := breaker.Load(nname); ok {
		return b.(*Breaker)
	}
	brk := NewBreakerWithOptions(&Options{
		Name: nname,
	})
	breaker.Store(nname, brk)
	return brk
}

func (c *Config) BreakerServer(name string) *Breaker {
	breaker := globalWatcher.breakers.Load().(*sync.Map)
	nname := fmt.Sprintf("%s.server.%s", c.namespace, name)
	if b, ok := breaker.Load(nname); ok {
		return b.(*Breaker)
	}
	brk := NewBreakerWithOptions(&Options{
		Name: nname,
	})
	breaker.Store(nname, brk)
	return brk
}
