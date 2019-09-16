package filter

import (
	"bytes"
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/discover"
	// "go_common_lib/http_client_pool"
	"go_common_lib/mytime"
	"time"

	"go_common_lib/go-json"
)

type FilterStatus struct {
}

func init() {
	BatchFilterMap["FilterStatus"] = FilterStatus{}
}
func (this FilterStatus) FilterInfos(req *data_type.Request) {
	if len(req.Livelist) == 0 {
		return
	}
	is_hit := false
	defer common.TimerV2("FilterStatus", &(req.Timer_log), time.Now(), &is_hit)
	var url bytes.Buffer
	// url.WriteString("http://")
	// ip := discover.GetUrl("phoenix", req)
	// url.WriteString(ip)
	url.WriteString(":18109/filter?rec_tab=")
	url.WriteString(req.Rec_tab)
	url.WriteString("&ids=")
	for idx, info := range req.Livelist {
		if idx != 0 {
			url.WriteString(",")
		}
		url.WriteString(info.LiveId)
	}
	result, err := discover.GetResult("phoenix", req, url.String(), 50)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	var rMap map[string]bool
	err = json.Unmarshal(result, &rMap)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	tempList := make([]data_type.LiveInfo, 0)
	for _, info := range req.Livelist {
		if _, isIn := rMap[info.LiveId]; isIn {
			tempList = append(tempList, info)
		}
	}
	req.Livelist = tempList
}
