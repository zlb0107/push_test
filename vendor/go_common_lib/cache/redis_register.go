package cache

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	Cache_redis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	//ri_array = append(ri_array, redis_pool.RedisInfo{&Cache_redis, "53cda00d4f19a97a5008aee662f04e7a.ali.codis.inkept.cn:6379", "hall_cont_rec:InkePassword", 500, 5000, 180, 100000000, 100000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Cache_redis, "r-2ze27ed0d7e1e504.redis.rds.aliyuncs.com:6379", "Us9g23fdWxi", 500, 5000, 180, 5000000, 5000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
