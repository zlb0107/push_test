package change

import (
	"go_common_lib/data_type"
)

type Change interface {
	ChangeData(*data_type.Request, chan string) int
}

var Change_map map[string]Change

//推荐架构可配置
var ChangeCfg_map map[string]func(moduleName string, params map[string]interface{}) Change

func init() {
	Change_map = make(map[string]Change)

	ChangeCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Change)
}
