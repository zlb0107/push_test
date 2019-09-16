package live_special_map

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	Baidu_redis      *redis.Pool
	Vip_redis        *redis.Pool
	FishRedis        *redis.Pool
	SexyRedis        *redis.Pool
	WLowRedis        *redis.Pool
	LLowRedis        *redis.Pool
	PerRecBlackRedis *redis.Pool
	MaybeLowRedis    *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	ri_array = append(ri_array, redis_pool.RedisInfo{&Baidu_redis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&MaybeLowRedis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Vip_redis, "r-2ze5ec929f2957d4.redis.rds.aliyuncs.com:6379", "Ne4w1Riy3", 50, 1000, 180, 500000000, 500000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&WLowRedis, "r-2ze5ec929f2957d4.redis.rds.aliyuncs.com:6379", "Ne4w1Riy3", 50, 1000, 180, 500000000, 500000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&FishRedis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&SexyRedis, "r-2ze5ec929f2957d4.redis.rds.aliyuncs.com:6379", "Ne4w1Riy3", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&LLowRedis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&PerRecBlackRedis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 50000000, 50000000})
	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
