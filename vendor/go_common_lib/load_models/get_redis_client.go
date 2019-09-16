package load_models

import (
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var pools sync.Map

func GetRedisClient(info *RedisInfo) redis.Conn {
	if strings.Index(info.RedisAddr, ":") == -1 {
		info.RedisAddr += ":6379"
	}

	instance, ok := pools.Load(info.RedisAddr)
	if ok {
		info.RedisClient = instance.(*redis.Pool)
	} else {
		redis_pool.RedisInit_v2(redis_pool.RedisInfo_v2{
			&(info.RedisClient),
			info.RedisAddr,
			info.RedisAuth,
			50,
			350,
			180,
			50 * time.Millisecond,
			20 * time.Millisecond,
			20 * time.Millisecond,
		})

		pools.Store(info.RedisAddr, info.RedisClient)
	}

	rc := info.RedisClient.Get()
	return rc
}
