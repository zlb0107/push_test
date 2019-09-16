// 低质量feed
package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/special_map/feed"
)

type FilterFeedLowQuality struct {
}

func init() {
	var rp FilterFeedLowQuality
	Filter_map["FilterFeedLowQuality"] = rp
	logs.Warn("in filter has rec init")
}

func (rp FilterFeedLowQuality) Filter_live(info *data_type.LiveInfo,
	request *data_type.Request) bool {
	return feed_special_map.LowQualityFeedidsController.Is_special_uid(info.LiveId)
}
