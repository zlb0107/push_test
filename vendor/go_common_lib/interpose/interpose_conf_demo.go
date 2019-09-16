package interpose

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"time"
)

type InterposeDemoCfg struct {
	ModuleName string
	Param      string
}

func init() {
	InterposeCfg_map["InterposeDemoCfg"] = NewInterposeDemoCfg
	logs.Warn("in InterposeDemoCfg init")
}

func NewInterposeDemoCfg(moduleName string, params map[string]interface{}) Interpose {
	var rp InterposeDemoCfg
	rp.ModuleName = moduleName
	rp.Param = params["param_b"].(string)

	return rp
}

func (rp InterposeDemoCfg) Run_interpose(request *data_type.Request) int {
	defer common.Timer("InterposeDemoCfg", &(request.Timer_log), time.Now())
	logs.Debug(rp.ModuleName, rp.Param)
	time.Sleep(1 * time.Second)

	return 0
}
