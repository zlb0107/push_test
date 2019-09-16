package live_special_map

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
	"go_common_lib/special_map"
	"time"
)

type StaticLiveIds struct {
	special_map.SpecialLiveIds
}

var StaticLiveIdsController StaticLiveIds

type StaticLiveIdsResp struct {
	DmError  int      `json:"dm_error"`
	ErrorMsg string   `json:"error_msg"`
	Data     []string `json:"data"`
}

func init() {
	go func() {
		for {
			StaticLiveIdsController.UpdateSimpleHttp("http://checkapi.busi.inke.cn/v1/live/static-list", getStaticLiveIds)
			time.Sleep(10 * time.Second)
		}
	}()
}
func getStaticLiveIds(staticLiveMap *map[string]bool, url string) {
	resp, err := http_client_pool.Get_url(url)
	if err != nil {
		logs.Error("error:", err)
		return
	}
	var stat_json StaticLiveIdsResp
	if err := json.Unmarshal(resp, &(stat_json)); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return
	}
	for _, liveId := range stat_json.Data {
		(*staticLiveMap)[liveId] = true
	}
}
