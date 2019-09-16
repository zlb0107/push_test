package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterDislikeFollow struct {
}

func init() {
	var rp FilterDislikeFollow
	Filter_map["FilterDislikeFollow"] = rp
	logs.Warn("in FilterDislikeFollow init")
}
func (rp FilterDislikeFollow) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Dislike_follow_list == nil {
		return false
	}
	if _, is_in := request.Dislike_follow_list[info.Uid]; is_in {
		return true
	}
	return false
}
