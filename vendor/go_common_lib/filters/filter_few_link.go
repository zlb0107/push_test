package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

// 多人直播间人数限制
// 多人视频直播间（6人间），当前处于保量或首页正常个性化推荐状态，开视频主播(含主持人)≤3个
//多人语音直播间（8人间），当前处于保量或首页正常个性化推荐状态，开语音主播(不含主持人)≤4个
type FilterFewLink struct {
}

func init() {
	var rp FilterFewLink
	Filter_map["FilterFewLink"] = rp
	logs.Warn("in FilterFewLink init")
}

func (rp FilterFewLink) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	subLiveType := living.Living_handler.GetSubLiveType(info.Uid)
	// 目标直播间为视频交友直播间和音频交友直播间
	if liveType != "friendlive" && (liveType != "audiolive" || subLiveType != "audiopal") && (liveType != "residentlive") {
		return false
	}
	// linkMikeNum 连麦的人数，不包含房主
	linkMikeNum := living.Living_handler.Get_live_linkMikeNum(info.Uid)

	// 没拿到不过滤
	if linkMikeNum < 0 {
		return false
	}
	if liveType == "friendlive" && linkMikeNum <= 2 {
		return true
	} else if liveType == "audiolive" && linkMikeNum < 3 {
		return true
	} else if (liveType == "residentlive") && linkMikeNum <= 2 {
		return true
	}
	return false
}
