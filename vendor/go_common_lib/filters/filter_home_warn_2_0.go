package filter

import (
	_ "sync"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	//all_filters "go_common_lib/filters"
	"go_common_lib/http_client_pool"
	_ "io/ioutil"
	_ "net/http"
	_ "strings"
	"time"

	"go_common_lib/go-json"
)

/*
	推荐首页过滤打散策略
	http://wiki.inkept.cn/pages/viewpage.action?pageId=31786340
	func:
		1.获取30分钟内的提示类过主播事件
		2.获取3小时内的警告类过滤主播事件
		3.过滤掉3小时内的警告类主播
*/

type FilterHomeWarn_2_0 struct {
}

type UidBadCaseStruct struct {
	Oper           string `json:"oper"`
	Last_oper_time int64  `json:"last_oper_time"`
}
type BadCaseStruct struct {
	Dm_errot       int                           `json:"dm_error"`
	Error_msg      string                        `json:"error_msg"`
	UidBadCaseInfo map[string][]UidBadCaseStruct `json:"data"`
}

const url = "http://checkapi.busi.inke.cn/v1/live/oper-data"
const timeout = 1000 //超时时间

const promptLimit = 30 * 60 * 1000      //提示类屏蔽时间界限 ms
const warningLimit = 3 * 60 * 60 * 1000 //警告类屏蔽时间界限 ms

/*保留相关信息*/
//30分钟内的信息
var BadCase30mInfoMap map[string][]UidBadCaseStruct

//3小时内的信息
var BadCase3hInfoMap map[string][]UidBadCaseStruct

func init() {
	var rp FilterHomeWarn_2_0
	Filter_map["FilterHomeWarn_2_0"] = rp
	logs.Warn("in FilterHomeWarn_2_0 init")
}

func (rp FilterHomeWarn_2_0) Init() {
	BadCase30mInfoMap = make(map[string][]UidBadCaseStruct, 0)
	BadCase3hInfoMap = make(map[string][]UidBadCaseStruct, 0)

	logs.Error("exec FilterHomeWarn_2_0 Init function")

	go func() {
		for {
			rp.refreshBadCase()
			time.Sleep(10 * time.Second)
		}
	}()
}

func (rp FilterHomeWarn_2_0) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	_, is_in := BadCase3hInfoMap[info.Uid]
	return is_in
}

func (rp FilterHomeWarn_2_0) refreshBadCase() {
	result, err := http_client_pool.Get_n_url(url, timeout)
	//logs.Debug("---result--:", string(result))
	if err != nil {
		logs.Warn("err:", err)
		return
	}
	var badCaseStruct BadCaseStruct
	err_parse := json.Unmarshal(result, &(badCaseStruct))
	if err_parse != nil {
		logs.Warn("err:", err_parse)
		return
	}
	//logs.Debug("--badCaseStruct--", badCaseStruct)
	if badCaseStruct.Dm_errot != 0 {
		logs.Warn("FilterHomeWarn_2_0 Dm_errot code return failed:", badCaseStruct)
		return
	}

	//30分钟内的信息
	var badCase30mInfoMap map[string][]UidBadCaseStruct = make(map[string][]UidBadCaseStruct, 0)
	//3小时内的信息
	var badCase3hInfoMap map[string][]UidBadCaseStruct = make(map[string][]UidBadCaseStruct, 0)

	curTime := time.Now().UnixNano() / 1000000
	for uid, uidBadCaseStruct := range badCaseStruct.UidBadCaseInfo {
		for _, uidBadCase := range uidBadCaseStruct {
			if uidBadCase.Oper == "提示" {
				if curTime-uidBadCase.Last_oper_time <= promptLimit {
					badCase30mInfoMap[uid] = append(BadCase30mInfoMap[uid], uidBadCase)
				}
			} else if uidBadCase.Oper == "警告" {
				if curTime-uidBadCase.Last_oper_time <= warningLimit {
					badCase3hInfoMap[uid] = append(badCase3hInfoMap[uid], uidBadCase)
				}
			} else {
				logs.Warn("FilterHomeWarn_2_0 Oper content is not identify:", uidBadCase.Oper)
			}
		}
	}
	/*
		if len(badCase30mInfoMap) != 0 { //后续增加过滤自己map内的超时时间
			BadCase30mInfoMap = badCase30mInfoMap
		}

		if len(badCase3hInfoMap) != 0 {
			BadCase3hInfoMap = badCase3hInfoMap
		}
	*/
	BadCase30mInfoMap = badCase30mInfoMap
	BadCase3hInfoMap = badCase3hInfoMap
}
