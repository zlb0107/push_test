package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

/*
Desc : 过滤电影
*/
type FilterFirm struct {
}

func init() {
	var rp FilterFirm
	Filter_map["FilterFirm"] = rp
	logs.Warn("init FilterFirm")
}

func (rp FilterFirm) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	subLiveType := living.Living_handler.GetSubLiveType(info.Uid)
	if liveType == "landscape" && subLiveType == "movie" {
		return true
	}

	return false
}
