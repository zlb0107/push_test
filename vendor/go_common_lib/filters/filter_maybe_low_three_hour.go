package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/live"
)

type FilterMaybeLowThreeHour struct {
}

func init() {
	var rp FilterMaybeLowThreeHour
	Filter_map["FilterMaybeLowThreeHour"] = rp
	logs.Warn("in FilterMaybeLowThreeHour init")
}
func (rp FilterMaybeLowThreeHour) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	if req.Page_idx > 0 {
		//只对首页生效
		return false
	}
	if live_special_map.MaybeLowThreeHourController.Is_special_uid(info.Uid) {
		return true
	}
	return false
}
