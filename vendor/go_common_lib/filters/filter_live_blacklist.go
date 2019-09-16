package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterLiveBlacklist struct {
}

func init() {
	var rp FilterLiveBlacklist
	Filter_map["FilterLiveBlacklist"] = rp
	logs.Warn("FilterLiveBlacklist init")
}
func (rp FilterLiveBlacklist) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Live_blacklist == nil {
		return false
	}
	if _, is_in := request.Live_blacklist[info.Uid]; is_in {
		//	logs.Error("has filter:", info.Uid)
		return true
	}
	return false
}
