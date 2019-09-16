package cache

import (
	logs "github.com/cihub/seelog"
	"github.com/coocood/freecache"
	//"runtime/debug"

	"go_common_lib/config"
)

/*
支持Get,Set,Del,不支持Mget
在app.conf配置文件中增加localcache::size配置，含义是缓冲的大小
fixme:
	1)key的最大长度不能超过65535
	2)key和value的最大长度不能超过(分配最大内存/256/4-24)
*/
var LocalCache *freecache.Cache

func init() {
	cacheSize, err := config.AppConfig.Int("localcache::size")
	if err != nil {
		cacheSize = 0 //freecache默认内存为512 * 1024
	}

	LocalCache = freecache.NewCache(cacheSize)
	go Cache_controlor.updateLocalCache()
}

func (this *CacheControl) updateLocalCache() {
	for {
		req := <-this.LocalChannel
		err := LocalCache.Set([]byte(req.Key), []byte(req.Value), int(req.Expire))
		if err != nil {
			logs.Error("put LocalCache err:", err)
		}
	}
}

func (this *CacheControl) UpdateLocalCache(key, value string, expire int64) {
	if len(this.LocalChannel) >= CHANNEL_SIZE/2 {
		//logs.Error("channel is busy:", len(this.Channel))
		return
	}

	req := Request{key, value, expire}
	this.LocalChannel <- req
}

func LocalCacheGet(key string) (string, error) {
	value, err := LocalCache.Get([]byte(key))
	return string(value), err
}

func LocalCacheDel(key string) bool {
	return LocalCache.Del([]byte(key))
}
