package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterFeedHasOut struct {
}

func init() {
	var rp FilterFeedHasOut
	Filter_map["FilterFeedHasOut"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterFeedHasOut) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.HasOutIds == nil {
		return false
	}
	if _, is_in := request.HasOutIds[info.LiveId]; is_in {
		return true
	}
	return false
}
