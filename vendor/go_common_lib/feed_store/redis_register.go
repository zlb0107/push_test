package store

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	StoreRedis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	//r-2ze4e8a00c585e84.redis.rds.aliyuncs.com
	//auth：Xw8jO0LyA
	//ri_array = append(ri_array, redis_pool.RedisInfo{&StoreRedis, "r-2ze4e8a00c585e84.redis.rds.aliyuncs.com:6379", "Xw8jO0LyA", 50, 1000, 180, 10000000000, 10000000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&StoreRedis, "suggest-feed-history.ali.codis.inkept.cn:6379", "feed-history:InkePassword", 50, 1000, 180, 20000000, 10000000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
