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
		3.过滤掉3小时内的提示类主播
*/

type FilterHomeFace_2_0 struct {
}

type UidBadCaseFaceStruct struct {
	Uid    string `json:"uid"`
	LiveId string `json:"live_id"`
}
type BadCaseFaceStruct struct {
	Dm_errot       int                    `json:"dm_error"`
	Error_msg      string                 `json:"error_msg"`
	UidBadCaseInfo []UidBadCaseFaceStruct `json:"data"`
}

const url_face = "http://checkapi.busi.inke.cn/v1/live/no-face"
const time_out_face = 1000

//30分钟内的信息
var BadCaseFaceMap map[string]UidBadCaseFaceStruct

func init() {
	var rp FilterHomeFace_2_0
	Filter_map["FilterHomeFace_2_0"] = rp
	logs.Warn("in FilterHomeFace_2_0 init")

	BadCaseFaceMap = make(map[string]UidBadCaseFaceStruct, 0)

	go func() {
		for {
			rp.refreshBadCase()
			time.Sleep(10 * time.Second)
		}
	}()
}
func (rp FilterHomeFace_2_0) Filter_live(info *data_type.LiveInfo, request *data_type.Request) bool {
	_, is_in := BadCaseFaceMap[info.Uid]
	return is_in
}

func (rp FilterHomeFace_2_0) refreshBadCase() {
	result, err := http_client_pool.Get_n_url(url_face, time_out_face)
	//logs.Debug("---result--:", string(result))
	if err != nil {
		logs.Warn("err:", err)
		return
	}
	var badCaseStruct BadCaseFaceStruct
	err_parse := json.Unmarshal(result, &(badCaseStruct))
	if err_parse != nil {
		logs.Warn("err:", err_parse)
		return
	}
	//logs.Debug("--badCaseStruct--", badCaseStruct)
	if badCaseStruct.Dm_errot != 0 {
		logs.Warn("FilterHomeFace_2_0 Dm_errot code return failed:", badCaseStruct)
		return
	}

	var badCaseFaceMap map[string]UidBadCaseFaceStruct = make(map[string]UidBadCaseFaceStruct, 0)

	for _, uidBadCaseFaceStruct := range badCaseStruct.UidBadCaseInfo {
		badCaseFaceMap[uidBadCaseFaceStruct.Uid] = uidBadCaseFaceStruct

	}
	/*
		if len(badCaseFaceMap) != 0 {
			BadCaseFaceMap = badCaseFaceMap
		}
	*/
	BadCaseFaceMap = badCaseFaceMap

}
