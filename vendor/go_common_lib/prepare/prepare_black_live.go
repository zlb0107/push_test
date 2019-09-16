package prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/live_black"
	"go_common_lib/mytime"
	_ "strconv"
	"time"
)

type BlackLiveList struct {
}

func init() {
	var rp BlackLiveList
	Prepare_map["BlackLiveList"] = rp
	logs.Warn("in BlackLiveList init")
}
func (rp BlackLiveList) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "BlackLiveList" }()
	defer common.Timer("BlackLiveList", &(request.Timer_log), time.Now())
	for idx, info := range request.Livelist {
		if live_black.Bm.Bad(info.Uid) {
			request.Livelist[idx].In_black_live = "1"
		} else {
			request.Livelist[idx].In_black_live = "0"
		}
	}
	return 0
}
