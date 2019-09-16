package filter

import (
	//"sync"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	//all_filters "go_common_lib/filters"
	//"go_common_lib/http_client_pool"
	//"time"
)

/*
	推荐首页过滤打散策略
	http://wiki.inkept.cn/pages/viewpage.action?pageId=31786340
	func:
	        1.依赖filter_home2_0_warn.go
		2.过滤掉30分钟内的提示类主播
*/
type FilterHomePrompt_2_0 struct {
}

func init() {
	var rp FilterHomePrompt_2_0
	Filter_map["FilterHomePrompt_2_0"] = rp
	logs.Warn("in FilterHomePrompt_2_0 init")
}

func (rp FilterHomePrompt_2_0) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	_, is_in := BadCase30mInfoMap[info.Uid]
	return is_in
}
