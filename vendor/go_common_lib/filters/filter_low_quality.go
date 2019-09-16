package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterLowQuality struct {
}

func init() {
	var rp FilterLowQuality
	Filter_map["FilterLowQuality"] = rp
	logs.Warn("in FilterLowQuality init")
}

func (rp FilterLowQuality) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	quality := living.Living_handler.GetQuality(info.Uid)
	if quality < -0.5 {
		return true
	}

	return false
}
