package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/live_black"
)

type FilterBlackList struct {
}

func init() {
	var rp FilterBlackList
	Filter_map["FilterBlackList"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterBlackList) Filter_live(info *data_type.LiveInfo, _ *data_type.Request) bool {
	if live_black.Bm.Bad(info.Uid) {
		return true
	}
	return false
}
