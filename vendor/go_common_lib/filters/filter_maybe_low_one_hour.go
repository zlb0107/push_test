package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/live"
)

type FilterMaybeLowOneHour struct {
}

func init() {
	var rp FilterMaybeLowOneHour
	Filter_map["FilterMaybeLowOneHour"] = rp
	logs.Warn("in FilterMaybeLowOneHour init")
}
func (rp FilterMaybeLowOneHour) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	if req.Page_idx > 0 {
		//只对首页生效
		return false
	}
	if live_special_map.MaybeLowOneHourController.Is_special_uid(info.Uid) {
		return true
	}
	return false
}
