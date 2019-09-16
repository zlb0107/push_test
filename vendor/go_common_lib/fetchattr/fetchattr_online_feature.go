package fetchattr

import (
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/snapshot"
)

type OnlineFetchattr struct {
}

func init() {
	var rp OnlineFetchattr
	Fetcher_map["OnlineFetchattr"] = rp
	logs.Warn("in OnlineFetchattr init")
}

func (of OnlineFetchattr) Get_attr(request *data_type.Request, ch chan string) int {
	defer common.Timer("OnlineFetchattr", &(request.Timer_log), time.Now())
	defer func() {
		ch <- "OnlineFetchattr"
	}()

	for _, info := range request.Livelist {
		snapshot.PutFeaturesToChannel(request.SnapshotVersion, info)
	}

	return 0
}
