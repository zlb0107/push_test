package prepare

import (
	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type Appearance struct {
}

func init() {
	var rp Appearance
	Prepare_map["Appearance"] = rp
	logs.Warn("in Appearance init")
}
func (rp Appearance) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "Appearance" }()
	//	defer common.Timer("Appearance", &(request.Timer_log), time.Now())
	rc := Appearance_redis.Get()
	is_hit := false
	defer common.TimerV2("Appearance", &(request.Timer_log), time.Now(), &is_hit)
	defer rc.Close()
	var keys []interface{}
	for _, info := range request.Livelist {
		keys = append(keys, "fl_"+info.Uid)
	}
	vs, err := cache.Get_redis_strings(keys, rc, &is_hit, 24*3600)
	//vs, err := redis.Values(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("get redis failed:", err)
		return -1
	}
	if len(vs) != len(request.Livelist) {
		logs.Error("vs.len != request.len ", len(vs), len(request.Livelist))
		return -1
	}
	for idx, v := range vs {
		if v == nil {
			request.Livelist[idx].Appearance = "-1"
			continue
		}
		request.Livelist[idx].Appearance = string(v.([]byte)) //1 好看 2 普通 3 丑
	}
	rp.Get_new_data(request)
	return 0
}
func (rp Appearance) Get_new_data(request *data_type.Request) int {
	rc := NewAppearanceRedis.Get()
	is_hit := false
	defer common.TimerV2("Appearance", &(request.Timer_log), time.Now(), &is_hit)
	defer rc.Close()
	var keys []interface{}
	for _, info := range request.Livelist {
		keys = append(keys, "newfl_"+info.Uid)
	}
	vs, err := cache.Get_redis_strings(keys, rc, &is_hit, 24*3600)
	//vs, err := redis.Values(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("get redis failed:", err)
		return -1
	}
	if len(vs) != len(request.Livelist) {
		logs.Error("vs.len != request.len ", len(vs), len(request.Livelist))
		return -1
	}
	for idx, v := range vs {
		if v == nil {
			continue
		}
		request.Livelist[idx].Appearance = string(v.([]byte)) //1 好看 2 普通 3 丑
	}
	return 0
}
