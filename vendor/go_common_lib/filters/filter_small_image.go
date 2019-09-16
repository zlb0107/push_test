package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/live"
)

type FilterImage struct {
}

func init() {
	var rp FilterImage
	Filter_map["FilterImage"] = rp
	logs.Warn("in filter low quality init")
}
func (rp FilterImage) Filter_live(info *data_type.LiveInfo, _ *data_type.Request) bool {
	return live_special_map.LiveSmallImageController.Is_special_uid(info.Uid)
}
