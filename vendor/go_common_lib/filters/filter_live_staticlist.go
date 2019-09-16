package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/live"
)

type FilterStaticLive struct {
}

func init() {
	var rp FilterStaticLive
	Filter_map["FilterStaticLive"] = rp
	logs.Warn("in filter FilterStaticLive init")
}
func (rp FilterStaticLive) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	return live_special_map.StaticLiveIdsController.Is_special_liveid(info.LiveId)
}
