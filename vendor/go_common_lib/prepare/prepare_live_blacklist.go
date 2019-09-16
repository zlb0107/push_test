package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	"go_common_lib/cache/lru"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/worker"

	"strings"
	"time"
)

var LiveBlacklistPrepareCache *lru.Cache

type LiveBlacklistPrepare struct {
}

func init() {
	var rp LiveBlacklistPrepare
	Prepare_map["LiveBlacklistPrepare"] = rp
	logs.Warn("in raw_prepare init")

	c, err := lru.NewCache(DEFAULT_CACHE_SIZE, "live_black_list")
	if err != nil {
		logs.Error("NewCache error: ", err)
		panic(err)
	}

	LiveBlacklistPrepareCache = c
	go worker.RunTask(LiveBlacklistPrepareCache, time.Second, time.Second, 1)
}

func (rp LiveBlacklistPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "LiveBlacklistPrepare" }()
	//defer common.Timer("LiveBlacklistPrepare", &(request.Timer_log), time.Now())
	is_hit := false
	defer common.TimerV2("LiveBlacklistPrepare", &(request.Timer_log), time.Now(), &is_hit)

	rc := User_blacklist_redis.Get()
	defer rc.Close()

	key := "zid_loathe_" + request.Uid
	strs, err := redis.String(LiveBlacklistPrepareCache.GetValue(
		key,
		60,
		func() (interface{}, error) { return rc.Do("get", key) },
		&is_hit,
	))
	if err != nil {
		logs.Error("get redis error: ", err, ", uid: ", request.Uid)
		return -1
	}

	//	strs, err := cache.Get_redis_string("zid_loathe_"+request.Uid, rc, &is_hit, 3600)
	//	if err != nil {
	//		return -1
	//	}

	//uid_loathe_153825592 10187042;46642836;40857899
	user_blacklist_map := make(map[string]bool)
	terms := strings.Split(strs, ";")
	for _, term := range terms {
		user_blacklist_map[term] = true
	}
	request.Live_blacklist = user_blacklist_map
	return 0
}
