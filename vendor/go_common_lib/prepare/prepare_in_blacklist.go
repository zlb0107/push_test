package prepare

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
)

const (
	InBlacklistUidPrefixKey  = "rec_b_u_"
	InBlacklistSmidPrefixKey = "rec_u_smid_"
)

// InBlacklistPrepare 判断用户是否是大数据黑名单用户
type InBlacklistPrepare struct {
}

func init() {
	var rp InBlacklistPrepare
	Prepare_map["InBlacklistPrepare"] = rp
	logs.Warn("in InBlacklistPrepare init")
}

func (rp InBlacklistPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "InBlacklistPrepare" }()
	is_hit := false
	defer common.TimerV2("InBlacklistPrepare", &(request.Timer_log), time.Now(), &is_hit)

	rc := BlacklistRedis.Get()
	defer rc.Close()

	var wg sync.WaitGroup
	// 只过滤smid就行
	//go prepareBlacklist(request, rc, InBlacklistUidPrefixKey+request.Uid)
	if request.Zid_str != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			prepareBlacklist(request, rc, InBlacklistSmidPrefixKey+request.Zid_str)
		}()
	}

	wg.Wait()
	return 0
}

func prepareBlacklist(request *data_type.Request, rc redis.Conn, key string) {
	ok, err := redis.Bool(rc.Do("EXISTS", key))
	if err != nil {
		logs.Error("is blacklist:", err)
	}

	if ok {
		request.IsBlacklist = true
	}
}
