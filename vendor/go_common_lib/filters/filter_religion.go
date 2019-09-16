// 宗教
package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterReligion struct {
}

func init() {
	var rp FilterReligion
	Filter_map["FilterReligion"] = rp
	logs.Warn("in FilterReligion init")
}

func (rp FilterReligion) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	return living.Living_handler.IsReligion(info.Uid)
}
