package lru

import (
	"sync"
	"time"

	"go_common_lib/open_falcon"

	glru "github.com/hashicorp/golang-lru"
)

type Cache struct {
	lruCache *glru.Cache

	sync.Mutex
	hit        int // 命中次数
	total      int // 总次数
	pluginName string
}

type Item struct {
	val    interface{}
	expire int64
}

func NewCache(size int, name string) (*Cache, error) {
	lruCache, err := glru.New(size)
	if err != nil {
		return nil, err
	}

	c := &Cache{
		lruCache:   lruCache,
		pluginName: name,
	}

	return c, nil
}

func (c *Cache) Working() error {
	now := time.Now()
	// 每小时上报一次
	if now.Second() == 0 && now.Minute() == 0 {
		c.Lock()
		total := c.total
		hit := c.hit

		c.total = 0
		c.hit = 0
		c.Unlock()

		tags := "plugin=" + c.pluginName
		rate := float64(hit) / float64(total) * 100
		return open_falcon.PostToOpenFalcon("event.cache.rate", 3600, rate, tags)
	}

	return nil
}

// 先从本地取
// 本地没有或者超过有效期，调用get方法取
func (c *Cache) GetValue(key interface{}, expire int64, get func() (interface{}, error), is_hit *bool) (interface{}, error) {
	c.Lock()
	c.total++
	c.Unlock()

	e, ok := c.lruCache.Get(key)
	if ok {
		i, ok := e.(*Item)
		if ok && i.expire > time.Now().Unix() {
			c.Lock()
			c.hit++
			c.Unlock()
			*is_hit = true
			return i.val, nil
		}
	}

	v, err := get()
	if err != nil {
		return nil, err
	}

	item := &Item{}
	item.val = v
	item.expire = time.Now().Unix() + expire
	c.lruCache.Add(key, item)
	return v, nil
}
