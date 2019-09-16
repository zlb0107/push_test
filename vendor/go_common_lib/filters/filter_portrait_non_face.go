// 封面为非真人头像
package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterPortraitNonFace struct {
}

func init() {
	var rp FilterPk
	Filter_map["FilterPortraitNonFace"] = rp
	logs.Warn("in FilterPortraitNonFace init")
}

func (rp FilterPortraitNonFace) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	has := living.Living_handler.PortraitHasFace(info.Uid)
	if !has {
		return true
	}

	return false
}
