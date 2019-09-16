package scatter

import (
	"go_common_lib/data_type"
)

type Scatter interface {
	Run_shuffle(*data_type.Request) int
	IsOk(data_type.Request, int) (int, bool)
}

var Scatter_map map[string]Scatter

//推荐架构可配置
var ScatterCfg_map map[string]func(moduleName string, params map[string]interface{}) Scatter

func init() {
	Scatter_map = make(map[string]Scatter)
	ScatterCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Scatter)
}
