package interpose

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"time"
)

type TopNewFeed struct {
}

func init() {
	var rp TopNewFeed
	Interpose_map["TopNewFeed"] = rp
	logs.Warn("in TopNewFeed init")
}
func (rp TopNewFeed) Run_interpose(request *data_type.Request) int {
	defer common.Timer("TopNewFeed", &(request.Timer_log), time.Now())
	top_list := make([]data_type.LiveInfo, 0)
	normal_list := make([]data_type.LiveInfo, 0)
	left_list := make([]data_type.LiveInfo, 0)
	top_num := 0
	for _, info := range request.Livelist {
		if top_num >= 2 {
			//已经选出，将其他顺序排好
			left_list = append(left_list, info)
		} else {
			if len(info.Token) > 8 && info.Token[:8] == "rec_12_1" {
				top_list = append(top_list, info)
				top_num += 1
			} else {
				normal_list = append(normal_list, info)
			}
		}
	}
	top_list = append(top_list, normal_list...)
	top_list = append(top_list, left_list...)
	request.Livelist = top_list

	return 0
}
