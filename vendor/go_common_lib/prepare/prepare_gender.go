package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	//"go_common_lib/cache"

	"go_common_lib/cache/lru"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/worker"

	"strings"
	"time"
)

var GenderPrepareCache *lru.Cache

type Gender struct {
}

func init() {
	var rp Gender
	Prepare_map["Gender"] = rp
	logs.Warn("in Gender init")

	c, err := lru.NewCache(DEFAULT_CACHE_SIZE, "gender")
	if err != nil {
		logs.Error("NewCache error: ", err)
		panic(err)
	}

	GenderPrepareCache = c
	go worker.RunTask(GenderPrepareCache, time.Second, time.Second, 1)
}

func (rp Gender) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "Gender" }()
	rc := Gender_redis.Get()
	defer rc.Close()

	is_hit := false
	defer common.TimerV2("Gender", &(request.Timer_log), time.Now(), &is_hit)

	key := "rec_u_b_" + request.Uid
	vs, err := redis.String(GenderPrepareCache.GetValue(
		key,
		3600,
		func() (interface{}, error) { return rc.Do("get", key) },
		&is_hit,
	))
	if err != nil {
		logs.Error("get redis error: ", err, ", uid: ", request.Uid)
		return -1
	}

	//	vs, err := cache.Get_redis_string(key, rc, &is_hit, 1*3600)
	//	//vs, err := redis.String(rc.Do("get", key))
	//	//redis-cli -h r-2ze5f71065637634368.redis.rds.aliyuncs.com -a vOByljzlvh26 get rec_u_b_238430817
	//	if err != nil {
	//		logs.Error("get redis failed:", err)
	//		return -1
	//	}

	//${gender}_${priority}|${age}_${priority}[|${other_info}_${other_priority}]
	if vs == "" {
		return 0
	}
	terms := strings.Split(vs, "|")
	if len(terms) < 1 {
		return 0
	}
	items := strings.Split(terms[0], "_")
	if len(items) != 2 {
		return 0
	}
	if items[0] != "1" && items[1] != "0" {
		return 0
	}
	request.Gender = items[0]
	return 0
}
