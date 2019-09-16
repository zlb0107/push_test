package feed_special_map

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	LowQualityFeedRedis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	ri_array = append(ri_array, redis_pool.RedisInfo{&LowQualityFeedRedis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 50000000, 50000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
