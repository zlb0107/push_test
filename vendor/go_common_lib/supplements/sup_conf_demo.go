package supplement

import (
	logs "github.com/cihub/seelog"
	//"github.com/garyburd/redigo/redis"
	//	"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"strings"
	"time"
)

type DemoSupCfg struct {
	ModuleName string
	RecList    string
}

func init() {
	SupplementCfg_map["DemoSupCfg"] = NewDemoSupCfg
	logs.Warn("in DemoSupCfg init")
}

func NewDemoSupCfg(moduleName string, params map[string]interface{}) Supplement {
	var rp DemoSupCfg
	rp.ModuleName = moduleName
	rp.RecList = params["rec_list"].(string)

	logs.Warn("moduleName:", moduleName, " params:", params)
	return rp
}

func (rp DemoSupCfg) Get_list(request *data_type.Request, c chan data_type.ChanStruct) int {
	var result data_type.ChanStruct
	result.Name = rp.ModuleName
	defer func() { c <- result }()
	if !request.Is_new {
		return 0
	}
	is_hit := false
	defer common.TimerV2(rp.ModuleName, &(request.Timer_log), time.Now(), &is_hit)

	recall_num := 0
	defer func() {
		new_log := *(request.Num_log) + " " + rp.ModuleName + "_num:" + strconv.Itoa(recall_num)
		request.Num_log = &new_log
	}()

	terms := strings.Split(rp.RecList, ";")
	//timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	num := 0
	for _, term := range terms {
		var info data_type.LiveInfo
		info.Uid = term
		info.LiveId = info.Uid
		//info.Token = "rec_7_10_1_0^" + request.Uid + "_" + timestamp + "_" + strconv.Itoa(num)
		info.Append_token = "^31"
		result.Livelist = append(result.Livelist, info)
		num++
		if num >= request.Count {
			break
		}
	}
	recall_num = len(result.Livelist)
	//time.Sleep(1 * time.Second)
	for _, info := range result.Livelist {
		logs.Debug(" uid:", info.Uid, " token:", info.Token)
	}
	return 0
}
