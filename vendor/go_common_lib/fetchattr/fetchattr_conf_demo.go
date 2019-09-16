package fetchattr

import (
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
)

type DemoCfgFetchattr struct {
	ModuleName string
	Param      string
}

func init() {
	FetcherCfg_map["DemoCfgFetchattr"] = NewDemoCfgFetchattr
	logs.Warn("in DemoCfgFetchattr init")
}

func NewDemoCfgFetchattr(moduleName string, params map[string]interface{}) AttrFetcher {
	var rp DemoCfgFetchattr
	rp.ModuleName = moduleName
	rp.Param = params["param_b"].(string)

	return rp
}

func (rp DemoCfgFetchattr) Get_attr(request *data_type.Request, ch chan string) int {
	defer common.Timer(rp.ModuleName, &(request.Timer_log), time.Now())
	defer func() {
		ch <- rp.ModuleName
	}()

	logs.Debug(rp.ModuleName, rp.Param)
	for idx, _ := range request.Livelist {
		request.Livelist[idx].Distance = 3.0
	}

	time.Sleep(1 * time.Second)
	return 0
}
