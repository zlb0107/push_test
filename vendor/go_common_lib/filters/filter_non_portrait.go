// 无头像
package filter

import (
	"strings"
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
	"go_common_lib/mytime"
)

type FilterNonPortrait struct {
}

func init() {
	var rp FilterNonPortrait
	Filter_map["FilterNonPortrait"] = rp
	logs.Warn("in FilterNonPortrait init")
}
func (rp FilterNonPortrait) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	defer common.Timer("FilterNonPortrait", &(req.Timer_log), time.Now())
	portrait := living.Living_handler.GetPortrait(info.Uid)
	if strings.HasSuffix(portrait, "MTUyODQyMzA0NTk2NiM2MjcjanBn.jpg") || portrait == "MzE3NDAxNDQ1Mzk1MTkx.jpg" || portrait == "NTIwNDcxNDUzNzE4ODAx.jpg" || portrait == "" {
		return true
	}
	return false
}
