package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterFeedHasRec struct {
}

func init() {
	var rp FilterFeedHasRec
	Filter_map["FilterFeedHasRec"] = rp
	logs.Warn("in filter has rec init")
}

func (rp FilterFeedHasRec) Filter_live(info *data_type.LiveInfo,
	request *data_type.Request) bool {
	if request.Has_rec_list == nil {
		return false
	}
	if _, is_in := request.Has_rec_list[info.LiveId]; is_in {
		return true
	}
	return false
}
