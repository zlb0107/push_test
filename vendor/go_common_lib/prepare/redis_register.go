package prepare

import (
	"github.com/garyburd/redigo/redis"
	"go_common_lib/my_redis_pool"
)

var (
	Less_love_redis        *redis.Pool
	Has_shown_redis        *redis.Pool
	Strong_filter_redis    *redis.Pool
	Appearance_redis       *redis.Pool
	NewAppearanceRedis     *redis.Pool
	User_blacklist_redis   *redis.Pool
	Hot_cluster_redis      *redis.Pool
	Has_rec_redis          *redis.Pool
	Quality_redis          *redis.Pool
	Face_redis             *redis.Pool
	Third_noportrait_redis *redis.Pool
	Gender_redis           *redis.Pool
	BigR_redis             *redis.Pool
	ClientGenderRedis      *redis.Pool
	BlacklistRedis         *redis.Pool
	NotActive_redis        *redis.Pool

	ShownRedis        *redis.Pool //feed show redis
	FilmHasShownRedis *redis.Pool
	BpcHasRecRedis    *redis.Pool
)

func init() {
	var ri_array []redis_pool.RedisInfo
	ri_array = append(ri_array, redis_pool.RedisInfo{&Less_love_redis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Has_shown_redis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Strong_filter_redis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Appearance_redis, "r-2zee82b066955c84359.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&NewAppearanceRedis, "r-2ze5ec929f2957d4.redis.rds.aliyuncs.com:6379", "Ne4w1Riy3", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&User_blacklist_redis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Hot_cluster_redis, "r-2ze9a3012577b554616.redis.rds.aliyuncs.com:6379", "m1GwbBzf6uvm", 50, 1000, 180, 50000000, 50000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Has_rec_redis, "53cda00d4f19a97a5008aee662f04e7a.ali.codis.inkept.cn:6379", "hall_cont_rec:InkePassword", 50, 1000, 180, 20000000, 200000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Quality_redis, "r-2zeb548c4c4a1784701.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 100000000, 100000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Face_redis, "r-2zee82b066955c84359.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 100000000, 100000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Third_noportrait_redis, "r-2zee82b066955c84359.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 100000000, 100000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&Gender_redis, "r-2ze5f71065637634368.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 2000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&BigR_redis, "r-2ze5331d4cea4d74464.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 2000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&ClientGenderRedis, "r-2ze5f71065637634368.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 100000000, 100000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&BlacklistRedis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 30000000, 30000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&NotActive_redis, "r-2ze37e214e4c9914.redis.rds.aliyuncs.com:6379", "Wx92Sjqa", 50, 1000, 180, 100000000, 100000000})

	ri_array = append(ri_array, redis_pool.RedisInfo{&ShownRedis, "r-2zec69cda56d50a4185.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})

	ri_array = append(ri_array, redis_pool.RedisInfo{&FilmHasShownRedis, "r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com:6379", "vOByljzlvh26", 50, 1000, 180, 20000000, 20000000})
	ri_array = append(ri_array, redis_pool.RedisInfo{&BpcHasRecRedis, "53cda00d4f19a97a5008aee662f04e7a.ali.codis.inkept.cn:6379", "hall_cont_rec:InkePassword", 50, 1000, 180, 20000000, 200000000})

	for idx, _ := range ri_array {
		redis_pool.RedisInit(ri_array[idx])
	}
}
