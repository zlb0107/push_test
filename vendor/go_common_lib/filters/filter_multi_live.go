package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterMultiLive struct {
}

func init() {
	var rp FilterMultiLive
	Filter_map["FilterMultiLive"] = rp
	logs.Warn("in FilterMultiLive init")
}

func (rp FilterMultiLive) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)

	if liveType == "friendlive" || liveType == "audiolive" {
		return true
	}

	return false
}
