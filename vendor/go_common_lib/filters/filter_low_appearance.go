package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
	"time"
)

type FilterLowAppearance struct {
}

func init() {
	var rp FilterLowAppearance
	Filter_map["FilterLowAppearance"] = rp
	logs.Warn("in FilterLowAppearance init")
}
func (rp FilterLowAppearance) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	hour := time.Now().Hour()
	if hour > 2 && hour < 6 {
		return false
	}
	if req.Page_idx > 1 {
		//大于2页的不做颜值限制
		return false
	}
	appearance := living.Living_handler.Get_appearance(info.Uid)
	if appearance == "3" {
		return true
	}
	return false
}
