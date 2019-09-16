package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterUserBlacklist struct {
}

func init() {
	var rp FilterUserBlacklist
	Filter_map["FilterUserBlacklist"] = rp
	logs.Warn("FilterUserBlacklist init")
}
func (rp FilterUserBlacklist) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.User_blacklist == nil {
		return false
	}
	if _, is_in := request.User_blacklist[info.Uid]; is_in {
		//	logs.Error("has filter:", info.Uid)
		return true
	}
	return false
}
