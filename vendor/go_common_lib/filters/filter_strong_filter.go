package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterStrongFilter struct {
}

func init() {
	var rp FilterStrongFilter
	Filter_map["FilterStrongFilter"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterStrongFilter) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Strong_filter_uids == nil {
		return false
	}
	if _, is_in := request.Strong_filter_uids[info.Uid]; is_in {
		return true
	}
	return false
}
