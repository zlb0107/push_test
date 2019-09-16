// 无三方头像
package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type FilterThirdNonPortrait struct {
}

func init() {
	var rp FilterThirdNonPortrait
	Filter_map["FilterThirdNonPortrait"] = rp
	logs.Warn("in FilterThirdNonPortrait init")
}

func (rp FilterThirdNonPortrait) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	portrait := living.Living_handler.GetThirdNonPortrait(info.Uid)
	if portrait == "0" {
		return true
	}

	return false
}
