package supplement

import (
	"strconv"
	"time"

	logs "github.com/cihub/seelog"

	"go_common_lib/data_type"
	"go_common_lib/discover"
	"go_common_lib/go-json"
	"go_common_lib/mytime"
)

type ZeusSup struct {
}

func init() {
	var rp ZeusSup
	Supplement_map["ZeusSup"] = rp
	logs.Warn("in ZeusSup init")
}
func (rp ZeusSup) Get_list(request *data_type.Request, c chan data_type.ChanStruct) int {
	var result data_type.ChanStruct
	result.Name = "ZeusSup"
	result.Livelist = make([]data_type.LiveInfo, 0)
	defer func() { c <- result }()
	is_hit := false
	defer common.TimerV2("ZeusSup", &(request.Timer_log), time.Now(), &is_hit)
	recall_num := 0
	filter_num := 0
	defer func() {
		new_log := *(request.Num_log) + " " + "ZeusSup_" + "num:" + strconv.Itoa(recall_num) + " ThunderFilter_num:" + strconv.Itoa(filter_num)
		request.Num_log = &new_log
	}()

	urlPostfix := ":18091/zeus?uid=" + request.Uid + "&count=" + strconv.Itoa(request.Count) + "&logid=" + request.Logid
	body, err := discover.GetResult("zeus", request, urlPostfix, 50)
	if err != nil {
		logs.Error("error:", err)
		return -1
	}
	Livelist := make([]data_type.LiveInfo, 0)
	if err := json.Unmarshal([]byte(body), &(Livelist)); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return -1
	}
	result.Livelist = Livelist
	recall_num = len(result.Livelist)
	return 0
}
