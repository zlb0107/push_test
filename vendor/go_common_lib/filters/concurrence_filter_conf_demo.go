package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"time"
)

type ConcurrenceFilterDemoCfg struct {
	ModuleName string
	Param      string
}

func init() {
	BatchFilterCfg_map["ConcurrenceFilterDemoCfg"] = NewConcurrenceFilterDemoCfg
}

func NewConcurrenceFilterDemoCfg(moduleName string, params map[string]interface{}) BatchFilter {
	var rp ConcurrenceFilterDemoCfg
	rp.ModuleName = moduleName

	rp.Param = params["param_b"].(string)

	return rp
}

func (rp ConcurrenceFilterDemoCfg) FilterInfos(req *data_type.Request) {
	if len(req.Livelist) == 0 {
		return
	}
	is_hit := false
	defer common.TimerV2("FileterStatus", &(req.Timer_log), time.Now(), &is_hit)
	filter_num := 0
	defer func() {
		new_log := *(req.Num_log) + " FilterStatus_num:" + strconv.Itoa(filter_num)
		req.Num_log = &new_log
	}()

	logs.Debug("param:", rp.Param)
	//time.Sleep(1 * time.Second)
	return
}
