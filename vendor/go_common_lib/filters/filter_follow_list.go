package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterFollowList struct {
}

func init() {
	var rp FilterFollowList
	Filter_map["FilterFollowList"] = rp
	logs.Warn("in FilterFollowList init")
}
func (rp FilterFollowList) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Follow_list == nil {
		return false
	}
	if _, is_in := request.Follow_list[info.Uid]; is_in {
		return true
	}
	return false
}
