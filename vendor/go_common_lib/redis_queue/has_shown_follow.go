package redis_queue

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"strconv"
	"strings"
	"time"
)

type HasShownFollow struct {
	Info_chan chan ZidInfo
}
type ZidInfo struct {
	Uid       string
	Queue     int
	Timestamp string
	Zid_list  []string
}

var Has_shown_follow_handler HasShownFollow

const QUEUE_LEN = 100
const KEY_PREFIX = "f_h_s_"

func init() {
	Has_shown_follow_handler.Info_chan = make(chan ZidInfo, 10000)
	go Has_shown_follow_handler.Deal_chan()
}
func (this *HasShownFollow) Update(zids_str string, req *data_type.Request) {

	queue_len := get_redis(req)
	zids := strings.Split(zids_str, ",")
	var zid_info ZidInfo
	zid_info.Zid_list = make([]string, 0)
	zid_info.Uid = req.Uid
	zid_info.Queue = queue_len
	for _, zid := range zids {
		req.Has_shown_uids[zid] = true
		zid_info.Zid_list = append(zid_info.Zid_list, zid)
	}
	//1830272461   2028年的时间戳
	zid_info.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	this.Info_chan <- zid_info
	//	logs.Error("hehe")
	//	logs.Flush()
}
func get_redis(req *data_type.Request) int {
	if req.Has_shown_uids == nil {
		req.Has_shown_uids = make(map[string]bool)
	}
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	strs, err := redis.Strings(rc.Do("zrevrange", KEY_PREFIX+req.Uid, 0, QUEUE_LEN))
	if err != nil {
		logs.Error("err:", err)
		return -1
	}
	for _, str := range strs {
		//	logs.Error("idx:", idx, " str:", str)
		if str == "" {
			continue
		}
		req.Has_shown_uids[str] = true
	}
	return len(req.Has_shown_uids)
}
func (this *HasShownFollow) Deal_chan() {
	for {
		info := <-this.Info_chan
		//	logs.Error("here:", len(this.Info_chan), info.Queue, QUEUE_LEN)
		if info.Queue >= QUEUE_LEN {
			control_queue(info.Uid)
		}
		add_queue(info)
	}
	return
}
func control_queue(uid string) {
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	_, err := rc.Do("zremrangebyrank", KEY_PREFIX+uid, 0, QUEUE_LEN)
	if err != nil {
		logs.Error("err:", err)
		return
	}
}
func add_queue(info ZidInfo) {
	//logs.Error("here")
	//logs.Flush()
	rc := Has_shown_redis.Get()
	defer rc.Close()
	defer logs.Flush()
	if len(info.Zid_list) == 0 {
		return
	}
	var keys []interface{}
	keys = append(keys, KEY_PREFIX+info.Uid)
	for _, zid := range info.Zid_list {
		keys = append(keys, info.Timestamp)
		keys = append(keys, zid)
	}
	_, err := rc.Do("zadd", keys...)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	rc.Do("expire", KEY_PREFIX+info.Uid, 30*86400)
}
