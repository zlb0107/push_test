/*************************************************************************************************
* Copyright 2019 INKE inc.
* File: filter_real_online_num.go
* Author: mingqi.wu
* Create Date: Jun 21 , 2019
* Modify:
* Desc:
	过滤真实在线人数小于3的主播
*************************************************************************************************/
package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterRealOnlineNum struct {
}

func init() {
	var rp FilterRealOnlineNum
	Filter_map["FilterRealOnlineNum"] = rp
	logs.Warn("in FilterRealOnlineNum init")
}

func (rp FilterRealOnlineNum) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	isAudioLive := living.Living_handler.IsAudiolive(info.Uid)
	if isAudioLive {
		if living.Living_handler.GetRealOnlineNum(info.Uid) < 8 {
			return true
		}
	}

	return false
}
