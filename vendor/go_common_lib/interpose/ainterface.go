package interpose

import (
	"go_common_lib/data_type"
)

type Interpose interface {
	Run_interpose(*data_type.Request) int
}

var Interpose_map map[string]Interpose

//推荐架构可配置
var InterposeCfg_map map[string]func(moduleName string, params map[string]interface{}) Interpose

func init() {
	Interpose_map = make(map[string]Interpose)
	InterposeCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Interpose)
}
