package supplement

import (
	//	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type Supplement interface {
	Get_list(*data_type.Request, chan data_type.ChanStruct) int
}

var Supplement_map map[string]Supplement

//推荐架构可配置
var SupplementCfg_map map[string]func(moduleName string, params map[string]interface{}) Supplement

func init() {
	Supplement_map = make(map[string]Supplement)
	SupplementCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Supplement)
}

type SortScoreType []data_type.LiveInfo

func (c SortScoreType) Len() int {
	return len(c)
}
func (c SortScoreType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c SortScoreType) Less(i, j int) bool {
	return c[i].Score > c[j].Score
}
