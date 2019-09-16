package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	//"go_common_lib/pk"
)

type FilterDemoCfg struct {
	ModuleName string
	Param      string
}

func init() {
	FilterCfg_map["FilterDemoCfg"] = NewFilterDemoCfg
	logs.Warn("in FilterDemoCfg init")
}

func NewFilterDemoCfg(moduleName string, params map[string]interface{}) Filter {
	var rp FilterDemoCfg
	rp.ModuleName = moduleName
	rp.Param = params["param_b"].(string)

	return rp
}

func (rp FilterDemoCfg) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	logs.Debug("param:", rp.Param)
	return false
}
