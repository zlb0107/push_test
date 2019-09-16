// 视频交友推荐语
package fetchattr

import (
	"time"

	"go_common_lib/data_type"
	"go_common_lib/living"
	"go_common_lib/mytime"

	logs "github.com/cihub/seelog"
)

type FetchattrFirend struct {
}

func init() {
	var rp FetchattrFirend
	Fetcher_map["FetchattrFirend"] = rp
	logs.Warn("in FetchattrFirend init")
}

func (of FetchattrFirend) Get_attr(request *data_type.Request, ch chan string) int {
	defer common.Timer("FetchattrFirend", &(request.Timer_log), time.Now())
	defer func() {
		ch <- "FetchattrFirend"
	}()

	for idx, info := range request.Livelist {
		liveType := living.Living_handler.Get_live_type(info.Uid)
		if liveType == "friendlive" {
			request.Livelist[idx].RecReason = "视频交友"
		}
	}

	return 0
}
