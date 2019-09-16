package models

import (
	"go_common_lib/data_type"
	"go_common_lib/feature_prepare"
)

type Models interface {
	Predict(*data_type.Request, *map[string]feature_prepare.FeatureWrapperStruct, chan string) int
}

var ModelsMap map[string]Models

//推荐架构可配置
var ModelsCfg_map map[string]func(moduleName string, params map[string]interface{}) Models

func init() {
	ModelsMap = make(map[string]Models)

	ModelsCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Models)
}
