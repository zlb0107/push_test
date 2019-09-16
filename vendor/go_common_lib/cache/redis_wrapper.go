package cache

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

// strs, err := redis.String(rc.Do("get", "rec_filter_"+request.Uid))

func Get_redis_string(key string, rc redis.Conn, is_hit *bool, expire int64) (string, error) {
	//	return redis.String(rc.Do("get", key))
	value, err := Cache_controlor.Get(key)
	if err == nil {
		*is_hit = true
		return value, err
	} else {
		str, err := redis.String(rc.Do("get", key))
		if err == nil {
			//将结果存入缓存
			Cache_controlor.Update(key, str, expire)
		}
		return str, err
	}
}
func Get_redis_strings(keys []interface{}, rc redis.Conn, is_hit *bool, expire int64) ([]interface{}, error) {
	//	return redis.Values(rc.Do("mget", keys...))
	//需要对缓存已有的直接取，没有的再取原redis
	values, err := Cache_controlor.Mget(keys)
	if err != nil {
		logs.Error("failed:", err)
		//访问redis失败
		str, err := redis.Values(rc.Do("mget", keys...))
		if err == nil {
			//将结果存入缓存
			for idx, value := range str {
				if value == nil {
					continue
				}
				Cache_controlor.Update(string(keys[idx].(string)), string(value.([]byte)), expire)
			}
		}
		return str, err
	}
	need_re_find := make(map[interface{}]int)
	again_keys := make([]interface{}, 0)
	for idx, value := range values {
		if value == nil {
			need_re_find[keys[idx]] = idx
			again_keys = append(again_keys, keys[idx])
		}
	}
	if len(again_keys) == 0 {
		//缓存全部取到
		*is_hit = true
		return values, err
	}
	//需要再取一次
	again_strs, again_err := redis.Values(rc.Do("mget", again_keys...))
	if again_err != nil {
		//再次取失败，直接返回之前结果
		return values, err
	}
	for idx, value := range again_strs {
		if value == nil {
			continue
		}
		key := again_keys[idx]
		real_idx := need_re_find[key]
		values[real_idx] = value
		//重新写入
		Cache_controlor.Update(string(keys[real_idx].(string)), string(value.([]byte)), expire)
	}
	return values, err
}
