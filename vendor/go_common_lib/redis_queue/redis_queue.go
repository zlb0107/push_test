package redis_queue

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type RedisQueue struct {
	Info_chan chan ZidInfoStruct
	queueLen  int
	keyPrefix string
	chanLen   int
}
type ZidInfoStruct struct {
	Uid       string
	Queue     int
	Timestamp string
	Zid_list  []string
}

func (this *RedisQueue) Update(zids []string, uid string) {
	var zid_info ZidInfoStruct
	zid_info.Zid_list = make([]string, 0)
	zid_info.Uid = uid
	zid_info.Queue = len(zids)
	for _, zid := range zids {
		zid_info.Zid_list = append(zid_info.Zid_list, zid)
	}
	zid_info.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	this.Info_chan <- zid_info
}
func (this *RedisQueue) GetList(uid string) map[string]bool {
	Has_shown_uids := make(map[string]bool)
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	strs, err := redis.Strings(rc.Do("zrevrange", this.keyPrefix+uid, 0, this.queueLen))
	if err != nil {
		logs.Error("err:", err)
		return Has_shown_uids
	}
	for _, str := range strs {
		if str == "" {
			continue
		}
		Has_shown_uids[str] = true
	}
	return Has_shown_uids
}
func (this *RedisQueue) Deal_chan() {
	for {
		info := <-this.Info_chan
		if info.Queue > this.queueLen {
			this.control_queue(info.Uid, (info.Queue - this.queueLen))
		}
		this.add_queue(info)
	}
	return
}
func (this *RedisQueue) control_queue(uid string, dis int) {
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	_, err := rc.Do("zremrangebyrank", this.keyPrefix+uid, 0, dis)
	if err != nil {
		logs.Error("err:", err)
		return
	}
}
func (this *RedisQueue) add_queue(info ZidInfoStruct) {
	//logs.Error("here")
	//logs.Flush()
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	if len(info.Zid_list) == 0 {
		return
	}
	var keys []interface{}
	keys = append(keys, this.keyPrefix+info.Uid)
	for _, zid := range info.Zid_list {
		keys = append(keys, info.Timestamp)
		keys = append(keys, zid)
	}
	_, err := rc.Do("zadd", keys...)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	rc.Do("expire", this.keyPrefix+info.Uid, 30*86400)
}
