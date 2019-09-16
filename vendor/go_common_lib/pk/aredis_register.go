package pk

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	Pk_redis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	ri_array = append(ri_array, redis_pool.RedisInfo{&Pk_redis, "r-2zee82b066955c84359.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 10000000000, 10000000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
