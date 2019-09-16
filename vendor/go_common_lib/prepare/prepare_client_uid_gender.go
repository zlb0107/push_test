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

type ClientUidGender struct {
}

func init() {
	var rp ClientUidGender
	Prepare_map["ClientUidGender"] = rp
	logs.Warn("in ClientUidGender init")
}
func (rp ClientUidGender) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "ClientUidGender" }()
	rc := ClientGenderRedis.Get()
	is_hit := false
	defer common.TimerV2("ClientUidGender", &(request.Timer_log), time.Now(), &is_hit)
	defer rc.Close()

	if (request.Gender != "0") && (request.Gender != "1") {
		request.Gender = "0"
	}

	key := "rec_info_u_" + request.Uid
	vs, err := cache.Get_redis_string(key, rc, &is_hit, 1*3600)
	//vs, err := redis.String(rc.Do("get", key))
	//redis-cli -h r-2ze5f71065637634368.redis.rds.aliyuncs.com -a vOByljzlvh26 get rec_info_u_238430817

	if err != nil {
		logs.Error("get redis failed:", err)
		return -1
	}
	//${gender}_${priority}|${age}_${priority}[|${other_info}_${other_priority}]
	/*
		if vs == "" {
			return 0
		}
	*/
	terms := strings.Split(vs, ":")
	if len(terms) > 0 {
		if (terms[0] == "0") || (terms[0] == "1") {
			request.Gender = terms[0]
		}
	}

	return 0
}
