package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type Filter interface {
	Filter_live(*data_type.LiveInfo, *data_type.Request) bool
}

//批量过滤接口
type BatchFilter interface {
	FilterInfos(*data_type.Request)
}

var Filter_map map[string]Filter
var BatchFilterMap map[string]BatchFilter

//推荐架构可配置
var FilterCfg_map map[string]func(moduleName string, params map[string]interface{}) Filter
var BatchFilterCfg_map map[string]func(moduleName string, params map[string]interface{}) BatchFilter

func init() {
	Filter_map = make(map[string]Filter)
	BatchFilterMap = make(map[string]BatchFilter)

	FilterCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Filter)
	BatchFilterCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) BatchFilter)
	println("in filter interface init")
}
func Is_filter(info data_type.LiveInfo, request *data_type.Request, filters []string) bool {
	for _, interface_name := range filters {
		inter, ret := Filter_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		if inter.Filter_live(&info, request) {
			//logs.Error("filter:", interface_name, " uid:", info.Uid)
			return true
		}
	}
	return false
}
func IsFilter(info *data_type.LiveInfo, request *data_type.Request, filters []string) bool {
	for _, interface_name := range filters {
		inter, ret := Filter_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		if inter.Filter_live(info, request) {
			//logs.Error("filter:", interface_name, " uid:", info.Uid)
			return true
		}
	}
	return false
}
func IsFilterWithReason(info data_type.LiveInfo, request *data_type.Request, filters []string) (bool, string) {
	for _, interface_name := range filters {
		inter, ret := Filter_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		if inter.Filter_live(&info, request) {
			return true, interface_name
		}
	}
	return false, ""
}
