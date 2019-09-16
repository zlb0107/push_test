// 电商直播
package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterRetailLive struct {
}

func init() {
	var rp FilterRetailLive
	Filter_map["FilterRetailLive"] = rp
	logs.Warn("in FilterRetailLive init")
}

func (rp FilterRetailLive) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	return living.Living_handler.IsRetailLive(info.Uid)
}
