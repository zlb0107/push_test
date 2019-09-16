package prepare

import (
	//"feed_zeus/redis_pool"
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strings"
	"time"
)

type FeedHasShownPrepare struct {
}

func init() {
	var rp FeedHasShownPrepare
	Prepare_map["FeedHasShownPrepare"] = rp
	logs.Warn("in FeedHasShownPrepare init")
	//cache.Cache_controlor.Register("rec_filter_")
}
func (rp FeedHasShownPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "FeedHasShownPrepare" }()
	is_hit := false
	defer common.TimerV2("FeedHasShownPrepare", &(request.Timer_log), time.Now(), &is_hit)
	rc := ShownRedis.Get()
	defer rc.Close()
	key := "feed_f_" + request.Uid
	str, err := redis.String(rc.Do("get", key))
	if err != nil {
		logs.Error("err:", err, " key:", key)
		return -1
	}
	//fuid_fid
	terms := strings.Split(str, ";")
	tempMap := make(map[string]bool)
	for _, term := range terms {
		items := strings.Split(term, "_")
		if len(items) != 2 {
			continue
		}
		fid := items[1]
		tempMap[fid] = true
	}
	request.Has_shown_uids = tempMap
	return 0
}
