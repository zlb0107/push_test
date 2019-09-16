package filter

//import (
//	logs "github.com/cihub/seelog"
//	"go_common_lib/data_type"
//	"go_common_lib/live_white"
//)
//
//type FilterWhiteList struct {
//}
//
//func init() {
//	var rp FilterWhiteList
//	Filter_map["FilterWhiteList"] = rp
//	logs.Warn("in filter low quality init")
//}
//func (rp FilterWhiteList) Filter_live(info *data_type.LiveInfo, _ *data_type.Request) bool {
//	if !live_white.Bm.Is_white(info.Uid) {
//		return true
//	}
//	return false
//}
