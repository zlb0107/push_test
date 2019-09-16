package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterFeedHasShown struct {
}

func init() {
	var rp FilterFeedHasShown
	Filter_map["FilterFeedHasShown"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterFeedHasShown) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Has_shown_uids == nil {
		return false
	}
	if _, is_in := request.Has_shown_uids[info.LiveId]; is_in {
		return true
	}
	return false
}
