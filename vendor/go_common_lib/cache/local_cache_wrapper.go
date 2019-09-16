package cache

import (
	//logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// strs, err := redis.String(rc.Do("get", "rec_filter_"+request.Uid))

func GetLocalString(key string, rc redis.Conn, is_hit *bool, expire int64) (string, error) {
	//	return redis.String(rc.Do("get", key))
	value, err := LocalCacheGet(key)
	//logs.Debug("--------GetLocalString-----", key, "value:", value)
	if err == nil {
		*is_hit = true
		return value, err
	} else {
		str, err := redis.String(rc.Do("get", key))
		if err == nil {
			//将结果存入缓存
			Cache_controlor.UpdateLocalCache(key, str, expire)
		}
		return str, err
	}
}
