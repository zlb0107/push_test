package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

/*
Author:wumingqi
Create Time: 2019-08-08
Medify Time:

Desc : 过滤所有的电台
Other:
	峰哥
*/
type FilterAllAudioLive struct {
}

func init() {
	var rp FilterAllAudioLive
	Filter_map["FilterAllAudioLive"] = rp
	logs.Warn("init FilterAllAudioLive")
}
func (rp FilterAllAudioLive) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	//audiolive:个人电台,多人电台

	if liveType == "audiolive" {
		return true
	}
	return false
}
