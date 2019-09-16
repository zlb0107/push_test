package cache

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Request struct {
	Key    string
	Value  string
	Expire int64
}
type CacheControl struct {
	Is_init      bool
	Channel      chan Request
	LocalChannel chan Request
}

var Cache_controlor CacheControl

const CHANNEL_SIZE = 1000

func init() {
	//初始化整体结构
	Cache_controlor.init()
}

func (this *CacheControl) init() {
	this.Channel = make(chan Request, CHANNEL_SIZE)
	this.LocalChannel = make(chan Request, CHANNEL_SIZE)
	go this.update()
}
func (this *CacheControl) update() {
	time.Sleep(1 * time.Second) //等待redis连接初始化
	rc := Cache_redis.Get()
	defer rc.Close()
	for {
		req := <-this.Channel
		//没有过期时间的批量写入，索性退化成单条写入
		_, err := rc.Do("setex", req.Key, req.Expire, req.Value)
		if err != nil {
			rc.Close()
			rc = Cache_redis.Get()
			logs.Error("put redis err:", err)
		}
	}
}
func (this *CacheControl) Update(key, value string, expire int64) {
	if len(this.Channel) >= CHANNEL_SIZE/2 {
		//logs.Error("channel is busy:", len(this.Channel))
		return
	}
	req := Request{key, value, expire}
	this.Channel <- req
}

func (this *CacheControl) Get(key string) (string, error) {
	rc := Cache_redis.Get()
	defer rc.Close()
	return redis.String(rc.Do("get", key))
}
func (this *CacheControl) Mget(keys []interface{}) ([]interface{}, error) {
	rc := Cache_redis.Get()
	defer rc.Close()
	return redis.Values(rc.Do("mget", keys...))
}
