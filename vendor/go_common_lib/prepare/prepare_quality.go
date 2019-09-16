package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"time"
)

type QualityPrepare struct {
	func_name string
}

func init() {
	var rp QualityPrepare
	rp.func_name = "QualityPrepare"
	Prepare_map[rp.func_name] = rp
	logs.Warn("in raw_prepare init")
}
func (rp QualityPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "QualityPrepare" }()
	defer common.Timer(rp.func_name, &(request.Timer_log), time.Now())
	rc := Quality_redis.Get()
	defer rc.Close()
	var keys []interface{}
	if len(request.Livelist) < 1 {
		return 0
	}
	for _, live_info := range request.Livelist {
		//rc.Send("get", "rec_lq_"+live_info.Liveid)
		keys = append(keys, "rec_lq_"+live_info.LiveId)
	}
	vs, err := redis.Values(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("get redis failed:", err)
		return -1
	}
	if len(vs) != len(request.Livelist) {
		logs.Error("vs.len != request.len", len(vs), len(request.Livelist))
		return -1
	}
	for idx, v := range vs {
		if v == nil {
			continue
		}
		request.Livelist[idx].Quality, err = strconv.ParseFloat(string(v.([]byte)), 32)
		if err != nil {
			logs.Error("trans float failed:", err)
			continue
		}
	}
	return 0
}
