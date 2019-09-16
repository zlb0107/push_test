package prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	store "go_common_lib/feed_store"
	"go_common_lib/mytime"
	"time"
)

type HasOutPrepare struct {
}

func init() {
	var rp HasOutPrepare
	Prepare_map["HasOutPrepare"] = rp
	//prepare.Prepare_map["HasOutPrepare"] = rp
	logs.Warn("in HasOutPrepare init")
	//cache.Cache_controlor.Register("rec_filter_")
}
func (rp HasOutPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "HasOutPrepare" }()
	is_hit := false
	defer common.TimerV2("HasOutPrepare", &(request.Timer_log), time.Now(), &is_hit)
	request.Ids = store.GetShown(request.Uid)
	if request.HasOutIds == nil {
		request.HasOutIds = make(map[string]bool)
	}
	for _, id := range *(request.Ids) {
		request.HasOutIds[id] = true
	}
	return 0
}
