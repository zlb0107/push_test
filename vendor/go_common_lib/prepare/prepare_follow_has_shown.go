package prepare

import (
	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/redis_queue"
	_ "strconv"
	"time"
)

type FollowHasShown struct {
}

func init() {
	var rp FollowHasShown
	Prepare_map["FollowHasShown"] = rp
	logs.Warn("in FollowHasShown init")
}
func (rp FollowHasShown) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "FollowHasShown" }()
	defer common.Timer("FollowHasShown", &(request.Timer_log), time.Now())
	redis_queue.Has_shown_follow_handler.Update(request.Zid_str, request)
	return 0
}
