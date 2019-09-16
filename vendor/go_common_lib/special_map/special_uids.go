package special_map

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type SpecialUids struct {
	uids map[string]bool
}

func (this *SpecialUids) Update(redis_key, split_str string, SpecialRedis *redis.Pool) {
	rc := SpecialRedis.Get()
	defer rc.Close()
	str, err := redis.String(rc.Do("get", redis_key))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	items := strings.Split(str, split_str)
	temp_map := make(map[string]bool)
	for _, item := range items {
		temp_map[item] = true
	}
	this.uids = temp_map
	logs.Error("redis_key:", redis_key, " size:", len(temp_map))
}
func (this *SpecialUids) UpdateKeys(keys []interface{}, split_str string, SpecialRedis *redis.Pool) {
	if len(keys) == 0 {
		logs.Error("key len:", len(keys))
		return
	}
	rc := SpecialRedis.Get()
	defer rc.Close()
	strs, err := redis.Strings(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	temp_map := make(map[string]bool)
	for _, str := range strs {
		items := strings.Split(str, split_str)
		for _, item := range items {
			temp_map[item] = true
		}
	}
	this.uids = temp_map
	logs.Error("redis_key:", string(keys[0].(string)), " size:", len(temp_map))
}
func (this *SpecialUids) UpdateTwoLevel(redis_key, split_one_str, split_two_str string, SpecialRedis *redis.Pool) {
	rc := SpecialRedis.Get()
	defer rc.Close()
	str, err := redis.String(rc.Do("get", redis_key))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	terms := strings.Split(str, split_one_str)
	temp_map := make(map[string]bool)
	for _, term := range terms {
		items := strings.Split(term, split_two_str)
		temp_map[items[0]] = true
	}
	this.uids = temp_map
	logs.Error("redis_key:", redis_key, " size:", len(temp_map))
}
func (this *SpecialUids) UpdateTwoLevelKeys(keys []interface{}, split_one_str, split_two_str string, SpecialRedis *redis.Pool) {
	if len(keys) == 0 {
		logs.Error("key len:", len(keys))
		return
	}
	rc := SpecialRedis.Get()
	defer rc.Close()
	strs, err := redis.Strings(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	temp_map := make(map[string]bool)
	for _, str := range strs {
		terms1 := strings.Split(str, split_one_str)
		for _, term1 := range terms1 {
			items2 := strings.Split(term1, split_two_str)
			temp_map[items2[0]] = true
		}
	}
	this.uids = temp_map
	logs.Error("redis_key:", string(keys[0].(string)), " size:", len(temp_map))
}
func (this *SpecialUids) UpdateHttp(allUids []string, limitNum int, url string, singleGet func(*map[string]bool, string)) {
	//end
	tempMap := make(map[string]bool)
	uids_str := ""
	for idx, uid := range allUids {
		if (idx+1)%limitNum == 0 {
			singleGet(&tempMap, url+uids_str)
			uids_str = ""
		}
		if len(uids_str) == 0 {
			uids_str = uid
		} else {
			uids_str += "," + uid
		}
	}
	if uids_str != "" {
		singleGet(&tempMap, url+uids_str)
	}

	this.uids = tempMap
	logs.Error("url:", url, " size:", len(tempMap))
}
func (this *SpecialUids) UpdateSimpleHttp(url string, singleGet func(*map[string]bool, string)) {
	tempMap := make(map[string]bool)
	singleGet(&tempMap, url)
	this.uids = tempMap
	logs.Error("url:", url, " size:", len(tempMap))
}
func (this *SpecialUids) Is_special_uid(uid string) bool {
	return this.uids[uid]
}
