package prepare

import (
	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strings"
	"time"
)

type StrongFilterPrepare struct {
}

func init() {
	var rp StrongFilterPrepare
	Prepare_map["StrongFilterPrepare"] = rp
	logs.Warn("in raw_prepare init")
}
func (rp StrongFilterPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "StrongFilterPrepare" }()
	is_hit := false
	defer common.TimerV2("StrongFilterPrepare", &(request.Timer_log), time.Now(), &is_hit)
	rc := Has_shown_redis.Get()
	defer rc.Close()
	//strs, err := redis.String(rc.Do("get", "rec_strong_filter_"+request.Uid))
	strs, err := cache.Get_redis_string("rec_strong_filter_"+request.Uid, rc, &is_hit, 60)
	if err != nil {
		//logs.Error("get redis failed:", err)
		//		logs.Error("get redis failed:", err, " uid:", request.Uid)
		//		logs.Flush()
		return -1
	}
	//4836893:-0.5;88341318:-0.5;4296493:-0.5
	uid_shown_map := make(map[string]bool)
	terms := strings.Split(strs, ";")
	for _, term := range terms {
		uid_shown_map[term] = true
		//	logs.Error("has_shown:", term)
	}
	request.Strong_filter_uids = uid_shown_map
	return 0
}
