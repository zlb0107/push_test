package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FilterHasRec struct {
}

func init() {
	var rp FilterHasRec
	Filter_map["FilterHasRec"] = rp
	logs.Warn("in filter has rec init")
}

//翻页过滤之前的sessionid内的
func (rp FilterHasRec) Filter_live(info *data_type.LiveInfo,
	request *data_type.Request) bool {
	if request.Has_rec_list == nil {
		return false
	}
	key := gen_key(info)
	if _, is_in := request.Has_rec_list[key]; is_in {
		return true
	}
	return false
}

func gen_key(info *data_type.LiveInfo) string {
	return info.Uid
}
