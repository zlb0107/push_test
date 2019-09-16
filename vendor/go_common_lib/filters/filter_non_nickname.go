package filter

import (
	"go_common_lib/living"
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
)

type FilterNonNickName struct {
}

func init() {
	var rp FilterNonNickName
	Filter_map["FilterNonNickName"] = rp
	logs.Info("in FilterNonNickName init")
}
func (rp FilterNonNickName) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	defer common.Timer("FilterNonNickName", &(req.Timer_log), time.Now())
	nickName := living.Living_handler.GetNickName(info.Uid)
	// nickName = "" 表示未获取到昵称，未获取到时不过滤
	if nickName == "Inke"+info.Uid || nickName == "inke"+info.Uid {
		return true
	}
	return false
}
