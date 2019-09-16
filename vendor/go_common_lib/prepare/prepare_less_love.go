package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"strings"
	"time"
)

type LessLovePrepare struct {
}

func init() {
	var rp LessLovePrepare
	Prepare_map["LessLovePrepare"] = rp
	logs.Warn("in raw_prepare init")
}
func (rp LessLovePrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "LessLovePrepare" }()
	defer common.Timer("LessLovePrepare", &(request.Timer_log), time.Now())
	rc := Less_love_redis.Get()
	defer rc.Close()
	strs, err := redis.String(rc.Do("get", "live_less_love_u_"+request.Uid))
	if err != nil {
		//logs.Error("get redis failed:", err)
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
