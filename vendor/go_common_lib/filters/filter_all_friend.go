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

Desc :
	过滤所有的交友,包括多人交友直播间，常驻交友直播间,交友多人直播（默认四人）
Other:
	峰哥
*/
type FilterAllFriend struct {
}

var gFriendLtLst = []string{"friendlive", "residentlive", "multiplayer"}

func init() {
	var rp FilterAllFriend
	Filter_map["FilterAllFriend"] = rp
	logs.Warn("init FilterAllFriend")
}
func (rp FilterAllFriend) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)

	for _, friendType := range gFriendLtLst {
		if liveType == friendType {
			return true
		}
	}

	return false
}
