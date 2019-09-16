package sort

import (
	"go_common_lib/data_type"
)

type Sort interface {
	Run_sort(*data_type.Request) int
}

var SortMap map[string]Sort

//推荐架构可配置
var SortCfg_map map[string]func(moduleName string, params map[string]interface{}) Sort

func init() {
	SortMap = make(map[string]Sort)

	SortCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Sort)
	println("in sort interface init")
}
