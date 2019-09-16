package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterBrush struct {
}

func init() {
	var rp FilterBrush
	Filter_map["FilterBrush"] = rp
	logs.Warn("init FilterBrush")
}
func (rp FilterBrush) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	brushNum := living.Living_handler.IsBrushVideo(info.Uid)
	if liveType != "audiolive" && brushNum < 0 {
		return true
	}
	return false
}
