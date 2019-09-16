package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/live"
)

type FileterPerRecBlack struct {
}

func init() {
	var rp FileterPerRecBlack
	Filter_map["FileterPerRecBlack"] = rp
	logs.Warn("in FileterPerRecBlack init")
}
func (rp FileterPerRecBlack) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	if live_special_map.PerRecBlackController.Is_special_uid(info.Uid) {
		return true
	}
	return false
}
