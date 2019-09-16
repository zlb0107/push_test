package interpose

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/pk"
	_ "strconv"
	"time"
)

type HighAppearance struct {
}

func init() {
	var rp HighAppearance
	Interpose_map["HighAppearance"] = rp
	logs.Warn("in HighAppearance init")
}
func (rp HighAppearance) Run_interpose(request *data_type.Request) int {
	defer common.Timer("HighAppearance", &(request.Timer_log), time.Now())

	hour := time.Now().Hour()
	if hour >= 1 && hour < 7 {
		//1-7点不做限制
		return 0
	}
	if request.Page_idx != 0 {
		//不是第一页的话，不做置顶
		return 0
	}
	top_list := make([]data_type.LiveInfo, 0)
	normal_list := make([]data_type.LiveInfo, 0)
	left_list := make([]data_type.LiveInfo, 0)
	top_num := 0
	for _, info := range request.Livelist {
		if top_num >= 6 {
			//已经选出，将其他顺序排好
			left_list = append(left_list, info)
		} else {
			if info.Appearance == "1" && !pk.Pk_handler.Is_pk(info.Uid) {
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
