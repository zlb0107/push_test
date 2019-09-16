package prepare

import (
	"time"

	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	"go_common_lib/data_type"
	"go_common_lib/mytime"
)

type PrepareNotActive struct {
}

func init() {
	var rp PrepareNotActive
	Prepare_map["PrepareNotActive"] = rp
	logs.Warn("in PrepareNotActive init")
}
func (rp PrepareNotActive) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "PrepareNotActive" }()
	rc := NotActive_redis.Get()
	defer rc.Close()
	defer common.Timer("PrepareNotActive", &(request.Timer_log), time.Now())
	key := "active_" + request.Uid
	vs, err := redis.String(rc.Do("get", key))
	if err != nil {
		logs.Error("get NotActive_redis failed:", err)
		return -1
	}
	request.IsNotActive = vs == "0"
	return 0
}
