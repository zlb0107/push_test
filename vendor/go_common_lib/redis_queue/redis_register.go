package redis_queue

import (
	"github.com/garyburd/redigo/redis"

	"go_common_lib/config"
	"go_common_lib/my_redis_pool"
)

var (
	Has_shown_redis *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	//ri_array = append(ri_array, redis_pool.RedisInfo{&Has_shown_redis, "r-2ze37e214e4c9914.redis.rds.aliyuncs.com:6379", "Wx92Sjqa", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Has_shown_redis, config.AppConfig.String("shown_redis::addr"), config.AppConfig.String("shown_redis::auth"), 50, 1000, 180, 20000000, 20000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
