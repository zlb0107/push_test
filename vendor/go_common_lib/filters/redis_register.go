package filter

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	NicksPortraitsRedis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	ri_array = append(ri_array, redis_pool.RedisInfo{&NicksPortraitsRedis, "r-2ze5ec929f2957d4.redis.rds.aliyuncs.com:6379", "Ne4w1Riy3", 50, 1000, 180, 20000000, 20000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
