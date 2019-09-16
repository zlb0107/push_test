package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type Face struct {
}

func init() {
	var rp Face
	Prepare_map["Face"] = rp
	logs.Warn("in Face init")
}
func (rp Face) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "Face" }()
	defer common.Timer("Face", &(request.Timer_log), time.Now())
	rc := Face_redis.Get()
	defer rc.Close()
	var keys []interface{}
	for _, info := range request.Livelist {
		keys = append(keys, "hf_"+info.Uid)
	}
	vs, err := redis.Values(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("get redis failed:", err)
		for idx, _ := range request.Livelist {
			request.Livelist[idx].Face = "-1"
		}
		return -1
	}

	if len(vs) != len(request.Livelist) {
		logs.Error("vs.len != request.len ", len(vs), len(request.Livelist))
		return -1
	}

	for idx, v := range vs {
		if v == nil {
			request.Livelist[idx].Face = "-1"
			continue
		}
		request.Livelist[idx].Face = string(v.([]byte)) //0 常规人像 1 婴儿  2 非人像
		if request.Livelist[idx].Face == "" {           //get hf_461858525  真的会有空字符串出现
			request.Livelist[idx].Face = "-1"
		}
	}
	return 0
}
