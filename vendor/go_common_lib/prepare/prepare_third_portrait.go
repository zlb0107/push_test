package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type ThirdPortrait struct {
}

func init() {
	var rp ThirdPortrait
	Prepare_map["ThirdPortrait"] = rp
	logs.Warn("in raw_prepare init")
}
func (rp ThirdPortrait) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "ThirdPortrait" }()
	defer common.Timer("ThirdPortrait", &(request.Timer_log), time.Now())
	rc := Third_noportrait_redis.Get()
	defer rc.Close()
	var keys []interface{}
	for _, info := range request.Livelist {
		keys = append(keys, "portrait_"+info.Uid)
	}
	vs, err := redis.Values(rc.Do("mget", keys...))
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
			request.Livelist[idx].Third_portrait = "-1"
			continue
		}
		request.Livelist[idx].Third_portrait = string(v.([]byte))
	}
	return 0
}
