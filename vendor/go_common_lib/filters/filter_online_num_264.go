package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
	"time"
)

type FilterOnlineNum264 struct {
}

func init() {
	var rp FilterOnlineNum264
	Filter_map["FilterOnlineNum264"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterOnlineNum264) Filter_live(info *data_type.LiveInfo, _ *data_type.Request) bool {
	online_num := living.Living_handler.Get_online_num(info.Uid)
	//分时段限制在线人数
	limit := 500
	hour := time.Now().Hour()
	//1-7_0 12-24_1500
	if hour > 1 && hour < 7 {
		limit = 0
	} else if hour > 12 {
		limit = 1500
	}
	if online_num < limit && online_num >= 0 {
		return true
	}
	return false
}
