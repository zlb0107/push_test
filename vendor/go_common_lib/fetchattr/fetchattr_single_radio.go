// 个人电台推荐语
package fetchattr

import (
	"time"

	"go_common_lib/data_type"
	"go_common_lib/living"
	"go_common_lib/mytime"

	logs "github.com/cihub/seelog"
)

type FetchattrSingleRadio struct {
}

func init() {
	var rp FetchattrSingleRadio
	Fetcher_map["FetchattrSingleRadio"] = rp
	logs.Warn("in FetchattrSingleRadio init")
}

func (of FetchattrSingleRadio) Get_attr(request *data_type.Request, ch chan string) int {
	defer common.Timer("FetchattrSingleRadio", &(request.Timer_log), time.Now())
	defer func() {
		ch <- "FetchattrSingleRadio"
	}()

	for idx, info := range request.Livelist {
		liveType := living.Living_handler.Get_live_type(info.Uid)
		subLiveType := living.Living_handler.GetSubLiveType(info.Uid)
		if liveType == "audiolive" && subLiveType != "audiopal" {
			request.Livelist[idx].RecReason = "电台"
		}
	}

	return 0
}
