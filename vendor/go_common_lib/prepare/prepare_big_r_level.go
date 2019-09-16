package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	"go_common_lib/data_type"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type PrepareBigRLevel struct {
}

func init() {
	var rp PrepareBigRLevel
	Prepare_map["PrepareBigRLevel"] = rp
	logs.Warn("in PrepareBigRLevel init")
}
func (rp PrepareBigRLevel) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "PrepareBigRLevel" }()
	rc := BigR_redis.Get()
	defer rc.Close()
	defer common.Timer("PrepareBigRLevel", &(request.Timer_log), time.Now())
	key := "rec_r_u_" + request.Uid
	vs, err := redis.String(rc.Do("get", key))
	if err != nil {
		// 非大R用户没有返回
		//logs.Error("get BigR redis failed:", err)
		return -1
	}
	if vs == "" {
		return 0
	}

	request.BigRLevel = vs
	return 0
}
