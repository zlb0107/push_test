package prepare

import (
	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"strings"
	"time"
)

type LessLovePrepareV2 struct {
}

func init() {
	var rp LessLovePrepareV2
	Prepare_map["LessLovePrepareV2"] = rp
	logs.Warn("in raw_prepare init")
}
func (rp LessLovePrepareV2) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "LessLovePrepareV2" }()
	is_hit := false
	defer common.TimerV2("LessLovePrepareV2", &(request.Timer_log), time.Now(), &is_hit)
	rc := Less_love_redis.Get()
	defer rc.Close()
	//strs, err := redis.String(rc.Do("get", "less_love_a_u_"+request.Uid))
	strs, err := cache.Get_redis_string("less_love_a_u_"+request.Uid, rc, &is_hit, 3600)
	if err != nil {
		//logs.Error("get redis failed:", err)
		//	logs.Error("get redis failed:", err, " uid:", request.Uid)
		//	logs.Flush()
		return -1
	}
	//4836893:-0.5;88341318:-0.5;4296493:-0.5
	uid_score_map := make(map[string]float64)
	terms := strings.Split(strs, ";")
	for _, term := range terms {
		vs := strings.Split(term, ":")
		if len(vs) != 2 {
			logs.Error("length is not 2:", len(vs), " ", term, " ", request.Uid)
			continue
		}
		score, _ := strconv.ParseFloat(vs[1], 32)
		uid_score_map[vs[0]] = score
	}
	request.Less_love_uids = uid_score_map
	return 0
}
