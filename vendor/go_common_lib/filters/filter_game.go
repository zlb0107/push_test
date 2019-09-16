// 新用户过滤游戏
package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterGame struct {
}

func init() {
	var rp FilterGame
	Filter_map["FilterGame"] = rp
	logs.Warn("init FilterGame")
}
func (rp FilterGame) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	if req.Is_new && liveType == "game" {
		return true
	}
	return false
}
