package sort

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"time"
)

type SortDemoCfg struct {
	ModuleName string
	Param      string
}

func init() {
	SortCfg_map["SortDemoCfg"] = NewSortDemoCfg
	logs.Warn("in distance level sort init")
}

func NewSortDemoCfg(moduleName string, params map[string]interface{}) Sort {
	var rp SortDemoCfg
	rp.ModuleName = moduleName
	rp.Param = params["param_b"].(string)

	return rp
}

func (rp SortDemoCfg) Run_sort(request *data_type.Request) int {
	defer common.Timer("SortDemoCfg", &(request.Timer_log), time.Now())
	for i, _ := range request.Livelist {
		request.Livelist[i].Token += strconv.Itoa(i)
	}
	logs.Debug("moduleName:", rp.ModuleName, " param:", rp.Param)
	time.Sleep(1 * time.Second)
	return 0
}
