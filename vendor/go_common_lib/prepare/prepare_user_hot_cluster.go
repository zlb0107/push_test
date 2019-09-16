package prepare

import (
	logs "github.com/cihub/seelog"
	//"github.com/garyburd/redigo/redis"
	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strings"
	"time"
)

type UserHotClusterPrepare struct {
}

func init() {
	var rp UserHotClusterPrepare
	Prepare_map["UserHotClusterPrepare"] = rp
	logs.Warn("in raw_prepare init")
	//cache.Cache_controlor.Register("rec_filter_")
}
func (rp UserHotClusterPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "UserHotClusterPrepare" }()
	is_hit := false
	defer common.TimerV2("UserHotClusterPrepare", &(request.Timer_log), time.Now(), &is_hit)
	rc := Hot_cluster_redis.Get()
	defer rc.Close()
	strs, err := cache.Get_redis_string("lda_kmeans_label_test_"+request.Uid, rc, &is_hit, 3600)
	if err != nil {
		return -1
	}
	//6_20171214
	terms := strings.Split(strs, "_")
	if len(terms) != 2 {
		return -1
	}
	request.User_hot_cluster = terms[0]
	return 0
}
