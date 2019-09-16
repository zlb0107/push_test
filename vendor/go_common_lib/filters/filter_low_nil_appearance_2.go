package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
)

type FilterLowNilAppearance2 struct {
}

const (
	// DefaultPloyId 默认策略
	DefaultPloyTagId = "1671"
	// GlobalPloyTagId 全局分发
	GlobalPloyTagId = "1672"
	// BanTop4PloyTagId 前4不出
	BanTop4PloyTagId = "1673"
	// BanTop48PloyTagId 前48不出
	BanTop48PloyTagId = "1674"
)

//rec_7_30:钓鱼
var tokenPrefixList = []string{"rec_7_18", "rec_7_19", "rec_7_20", "rec_7_21", "rec_7_22", "rec_7_23", "rec_7_37", "rec_7_30"}

func init() {
	var rp FilterLowNilAppearance2
	Filter_map["FilterLowNilAppearance2"] = rp
	logs.Warn("in FilterLowNilAppearance2 init")
}

func (rp FilterLowNilAppearance2) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	// 电影主播、优质户外主播、签约主持人不做分发评级
	if len(info.Token) > 8 {
		for _, token := range tokenPrefixList {
			if info.Token[:8] == token {
				return false
			}
		}
	}

	if req.Page_idx > 1 {
		//大于2页的不做颜值限制
		return false
	}
	// 全局分发则不过滤
	if living.Living_handler.HasTag(info.Uid, GlobalPloyTagId) {
		return false
	} else if living.Living_handler.HasTag(info.Uid, BanTop48PloyTagId) {
		// 前48不出过滤
		return true
	}

	appearance := living.Living_handler.Get_appearance(info.Uid)
	if appearance != "1" && appearance != "2" && appearance != "4" {
		return true
	}
	return false
}
