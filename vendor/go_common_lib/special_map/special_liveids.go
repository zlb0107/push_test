package special_map

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type SpecialLiveIds struct {
	liveids map[string]bool
}

func (this *SpecialLiveIds) Update(redis_key, split_str string, SpecialRedis *redis.Pool) {
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
	this.liveids = temp_map
	logs.Error("redis_key:", redis_key, " size:", len(temp_map))
}
func (this *SpecialLiveIds) UpdateKeys(keys []interface{}, split_str string, SpecialRedis *redis.Pool) {
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
	this.liveids = temp_map
	logs.Error("redis_key:", string(keys[0].(string)), " size:", len(temp_map))
}
func (this *SpecialLiveIds) UpdateTwoLevel(redis_key, split_one_str, split_two_str string, SpecialRedis *redis.Pool) {
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
	this.liveids = temp_map
	logs.Error("redis_key:", redis_key, " size:", len(temp_map))
}
func (this *SpecialLiveIds) UpdateHttp(allLiveIds []string, limitNum int, url string, singleGet func(*map[string]bool, string)) {
	//end
	tempMap := make(map[string]bool)
	liveids_str := ""
	for idx, liveId := range allLiveIds {
		if (idx+1)%limitNum == 0 {
			singleGet(&tempMap, url+liveids_str)
			liveids_str = ""
		}
		if len(liveids_str) == 0 {
			liveids_str = liveId
		} else {
			liveids_str += "," + liveId
		}
	}
	if liveids_str != "" {
		singleGet(&tempMap, url+liveids_str)
	}

	this.liveids = tempMap
	logs.Error("url:", url, " size:", len(tempMap))
}
func (this *SpecialLiveIds) UpdateSimpleHttp(url string, singleGet func(*map[string]bool, string)) {
	tempMap := make(map[string]bool)
	singleGet(&tempMap, url)
	this.liveids = tempMap
	logs.Error("url:", url, " size:", len(tempMap))
}
func (this *SpecialLiveIds) Is_special_liveid(liveId string) bool {
	return this.liveids[liveId]
}
