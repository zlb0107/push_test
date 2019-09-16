package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterAudioLive struct {
}

func init() {
	var rp FilterAudioLive
	Filter_map["FilterAudioLive"] = rp
	logs.Warn("init FilterAudioLive")
}
func (rp FilterAudioLive) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	subLiveType := living.Living_handler.GetSubLiveType(info.Uid)
	if liveType == "audiolive" {
		if subLiveType == "audiopal" {
			//多人电台不做过滤
			return false
		}
		return true
	}
	return false
}
