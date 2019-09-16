package prepare

import (
	logs "github.com/cihub/seelog"
	//"github.com/garyburd/redigo/redis"
	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"strings"
	"time"
)

type MultiUidGender struct {
}

func init() {
	var rp MultiUidGender
	Prepare_map["MultiUidGender"] = rp
	logs.Warn("in MultiUidGender init")
}
func (rp MultiUidGender) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "MultiUidGender" }()
	rc := ClientGenderRedis.Get()
	is_hit := false
	defer common.TimerV2("MultiUidGender", &(request.Timer_log), time.Now(), &is_hit)
	defer rc.Close()

	if len(request.Livelist) == 0 {
		return 1
	}

	var keys []interface{}
	for _, info := range request.Livelist {
		keys = append(keys, "rec_info_u_"+info.Uid)
	}
	//vs, err := redis.Values(rc.Do("mget", keys...))
	//redis-cli -h r-2ze5f71065637634368.redis.rds.aliyuncs.com -a vOByljzlvh26 get rec_info_u_238430817
	vs, err := cache.Get_redis_strings(keys, rc, &is_hit, 1*3600)
	if err != nil {
		logs.Error("get redis failed:", err)
		for idx, _ := range request.Livelist {
			request.Livelist[idx].Gender = "0"
		}
		return -1
	}
	if len(vs) != len(request.Livelist) {
		logs.Error("vs.len != request.len ", len(vs), len(request.Livelist))
		for idx, _ := range request.Livelist {
			request.Livelist[idx].Gender = "0"
		}
		return -1
	}

	for idx, v := range vs {
		//${gender}_${priority}|${age}_${priority}[|${other_info}_${other_priority}]
		if v == nil {
			request.Livelist[idx].Gender = "0"
			continue
		}
		terms := strings.Split(string(v.([]byte)), ":")
		if len(terms) < 1 {
			request.Livelist[idx].Gender = "0"
			continue
		}

		if (terms[0] == "0") || (terms[0] == "1") {
			request.Livelist[idx].Gender = terms[0]
		} else {
			request.Livelist[idx].Gender = "0"
		}
	}

	return 0
}
