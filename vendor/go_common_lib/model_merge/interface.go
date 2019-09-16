package model_merge

import (
	"go_common_lib/data_type"
)

type ModelMerge interface {
	Merge(*data_type.Request) int
}

var ModelMergeMap map[string]ModelMerge

//推荐架构可配置
var ModelMergeCfg_map map[string]func(moduleName string, params map[string]interface{}) ModelMerge

func init() {
	ModelMergeMap = make(map[string]ModelMerge)

	ModelMergeCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) ModelMerge)
}
