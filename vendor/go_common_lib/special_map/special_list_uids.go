package special_map

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type SpecialListUids struct {
	Uids   []string
	RWLock sync.RWMutex
	Size   int
}

func (this *SpecialListUids) GetAllUids() []string {
	this.RWLock.RLock()
	defer this.RWLock.RUnlock()
	return this.Uids
}

func (this *SpecialListUids) GetRandUids(num int) (uids []string, size int) {
	this.RWLock.RLock()
	defer this.RWLock.RUnlock()
	size = this.Size
	rand.Seed(time.Now().UnixNano())
	baseIdx := rand.Intn(this.Size)
	if this.Size-baseIdx >= num {
		uids = this.Uids[baseIdx : baseIdx+num]
	} else {
		lastNum := num + baseIdx - this.Size
		uids = this.Uids[baseIdx:]
		if lastNum > 1 {
			uids = append(uids, this.Uids[0:lastNum]...)
		} else {
			uids = append(uids, this.Uids[0])
		}

	}

	return
}

func (this *SpecialListUids) Update(redis_key, split_str string, SpecialRedis *redis.Pool) {
	rc := SpecialRedis.Get()
	defer rc.Close()
	str, err := redis.String(rc.Do("get", redis_key))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	items := strings.Split(str, split_str)
	this.RWLock.Lock()
	this.Uids = items
	this.Size = len(items)
	this.RWLock.Unlock()
	logs.Error("redis_key:", redis_key, " size:", len(items))
}
func (this *SpecialListUids) UpdateKeys(keys []interface{}, split_str string, SpecialRedis *redis.Pool) {
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
	var tmp_lst []string
	for _, str := range strs {
		items := strings.Split(str, split_str)
		if len(items) != 0 {
			tmp_lst = append(tmp_lst, items...)
		}
	}
	this.RWLock.Lock()
	this.Uids = tmp_lst
	this.Size = len(tmp_lst)
	this.RWLock.Unlock()
	logs.Error("redis_key:", string(keys[0].(string)), " size:", len(tmp_lst))
}
func (this *SpecialListUids) UpdateTwoLevel(redis_key, split_one_str, split_two_str string, SpecialRedis *redis.Pool) {
	rc := SpecialRedis.Get()
	defer rc.Close()
	str, err := redis.String(rc.Do("get", redis_key))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	terms := strings.Split(str, split_one_str)
	var tmp_lst []string
	for _, term := range terms {
		items := strings.Split(term, split_two_str)
		if (len(items) > 0) && (items[0] != "") {
			tmp_lst = append(tmp_lst, items[0])
		}
	}
	this.RWLock.Lock()
	this.Uids = tmp_lst
	this.Size = len(tmp_lst)
	this.RWLock.Unlock()
	logs.Error("redis_key:", redis_key, " size:", len(tmp_lst))
}
