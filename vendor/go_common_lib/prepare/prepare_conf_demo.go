package prepare

import (
	//"fmt"
	logs "github.com/cihub/seelog"
	//"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	// "strconv"
	"time"
)

type DemoPrepareCfg struct {
	ModuleName string
	Param_a    int
	Param_b    string
}

func init() {
	PrepareCfg_map["DemoPrepareCfg"] = NewDemoPrepareCfg
	logs.Warn("in DemoPrepareCfg init")
}

//如果初始化失败，直接panic
//func NewDemoPrepareCfg(moduleName string, params map[string]interface{}) (rp DemoPrepareCfg) {
func NewDemoPrepareCfg(moduleName string, params map[string]interface{}) Prepare {
	var rp DemoPrepareCfg
	rp.ModuleName = moduleName
	/*
		if len(params) != 3 { // template
			panic("NewDemoPrepareCfg param's len is error, moduleName:" + moduleName)
		}
	*/
	//logs.Warn("moduleName:",moduleName, " template:", params["template"], " params:",params)
	logs.Warn("moduleName:", moduleName, " params:", params, "  :", params["param_b"])

	rp.Param_a = int(params["param_a"].(float64))

	rp.Param_b = params["param_b"].(string)

	//rp.Param_b = string(params["param_b"].(string))
	/*
			if v, isIn := params["param_a"]; isIn {
				if f, isIn := v.(float64); isIn {
					rp.Param_a = int(f)
				} else {
					panic("conveer failed")
				}
			}


		if v, isIn := params["param_b"]; isIn {
			if val, isOk := v.(string); isOk {
				rp.Param_b = val
			} else {
				panic("param_b convert failed")
			}
		} else {
			panic("param_b failed")
		}
	*/

	return rp
}

func (rp DemoPrepareCfg) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- rp.ModuleName }()
	is_hit := false
	defer common.TimerV2(rp.ModuleName, &(request.Timer_log), time.Now(), &is_hit)
	for idx, _ := range request.Livelist {
		//request.Livelist[idx].Gender = rp.Param_b
		request.Livelist[idx].Portrait = rp.Param_b
	}

	logs.Warn("moduleName:", rp.ModuleName, " params:", rp.Param_a, "  :", rp.Param_b)

	/*
		for _, info := range request.Livelist {
			logs.Debug("result:", info.Uid, " token:", info.Token)
		}
	*/
	//time.Sleep(1 * time.Second)

	return 0
}
