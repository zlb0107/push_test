package change

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type ChangeDemoCfg struct {
	ModuleName string
	Param      string
}

func init() {
	ChangeCfg_map["ChangeDemoCfg"] = NewChangeDemoCfg
	logs.Warn("ChangeDemoCfg init")
}

func NewChangeDemoCfg(moduleName string, params map[string]interface{}) Change {
	var rp ChangeDemoCfg
	rp.ModuleName = moduleName
	rp.Param = params["param_b"].(string)

	return rp
}

func (rp ChangeDemoCfg) ChangeData(request *data_type.Request, c chan string) int {
	defer func() { c <- rp.ModuleName }()
	logs.Debug(rp.ModuleName, " param:", rp.Param)
	return 0
}
