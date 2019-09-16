package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterDislikeHall struct {
}

func init() {
	var rp FilterDislikeHall
	Filter_map["FilterDislikeHall"] = rp
	logs.Warn("in FilterDislikeHall init")
}
func (rp FilterDislikeHall) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	if request.Dislike_hall_list == nil {
		return false
	}
	if _, is_in := request.Dislike_hall_list[info.Uid]; is_in {
		return true
	}
	return false
}
