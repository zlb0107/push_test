package filter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

/*
Desc : 过滤多人电台
*/
type FilterAudiopal struct {
}

func init() {
	var rp FilterAudiopal
	Filter_map["FilterAudiopal"] = rp
	logs.Warn("init FilterAudiopal")
}

func (rp FilterAudiopal) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	subLiveType := living.Living_handler.GetSubLiveType(info.Uid)
	if liveType == "audiolive" && subLiveType == "audiopal" {
		return true
	}

	return false
}
