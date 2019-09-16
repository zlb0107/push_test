package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterLowNilAppearance struct {
}

func init() {
	var rp FilterLowNilAppearance
	Filter_map["FilterLowNilAppearance"] = rp
	logs.Warn("in FilterLowNilAppearance init")
}
func (rp FilterLowNilAppearance) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	if req.Page_idx > 1 {
		//大于2页的不做颜值限制
		return false
	}
	appearance := living.Living_handler.Get_appearance(info.Uid)
	if appearance != "1" && appearance != "2" && appearance != "4" {
		return true
	}
	return false
}
