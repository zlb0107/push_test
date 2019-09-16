package prepare

import (
	"go_common_lib/data_type"
)

const (
	DEFAULT_CACHE_SIZE = 100000
)

type Prepare interface {
	Get_data(*data_type.Request, chan string) int
}

var Prepare_map map[string]Prepare

//推荐架构可配置
var PrepareCfg_map map[string]func(moduleName string, params map[string]interface{}) Prepare

func init() {
	Prepare_map = make(map[string]Prepare)
	PrepareCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Prepare)
}
