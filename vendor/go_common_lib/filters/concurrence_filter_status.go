package filter

import (
	"bytes"
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/discover"
	// "go_common_lib/http_client_pool"
	"go_common_lib/load_confs"
	"go_common_lib/mytime"
	"strconv"
	"time"

	"go_common_lib/go-json"
)

type ConcurrenceFilterStatus struct {
}

func init() {
	BatchFilterMap["ConcurrenceFilterStatus"] = ConcurrenceFilterStatus{}
}
func (this ConcurrenceFilterStatus) FilterInfos(req *data_type.Request) {
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

	var url bytes.Buffer
	// url.WriteString("http://")
	// ip := discover.GetUrl("phoenix", req)
	// url.WriteString(ip)
	//url.WriteString(":18109/filter_get_leader?rec_tab=")
	url.WriteString(":18109/filter_get_uid_leader?rec_tab=")
	url.WriteString(req.Rec_tab)
	url.WriteString("&uid=")
	url.WriteString(req.Uid)
	url.WriteString("&logids=")
	url.WriteString(req.Logid)
	url.WriteString("&ids=")
	concurrenceNum := 100
	start := 0
	end := 0
	routineNum := func() int {
		if len(req.Livelist)%concurrenceNum == 0 {
			return len(req.Livelist) / concurrenceNum
		}
		return (len(req.Livelist) / concurrenceNum) + 1
	}()
	routineChan := make(chan *map[string]data_type.PhoenixResp, routineNum)
	hasInMap := make(map[string]bool)
	for i := 0; i < routineNum; i++ {
		start = i * concurrenceNum
		end = (i + 1) * concurrenceNum
		if end > len(req.Livelist) {
			end = len(req.Livelist)
		}

		//start , end 左闭右开
		go ConcurenceGet(req, url.String(), start, end, routineChan)
	}
	tempMap := make(map[string]data_type.PhoenixResp)
	for i := 0; i < routineNum; i++ {
		arrPtr := <-routineChan
		if arrPtr == nil {
			continue
		}
		for feedid, uidAndLeader := range *arrPtr {
			tempMap[feedid] = uidAndLeader
		}
	}
	tempList := make([]data_type.LiveInfo, 0)
	//exp_strategy := load_confs.ExpMap[req.Mylogid]
	exp_strategy, ok := req.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return
	}
	for _, info := range req.Livelist {
		//是否被过滤掉
		uidAndLeader, isIn := tempMap[info.LiveId]
		//leader, isIn := tempMap[info.LiveId].LeaderId
		if !isIn {
			filter_num += 1
			continue
		}
		leader := uidAndLeader.LeaderId
		//是否属于需要过滤的主播
		if info.LiveId != leader {
			info.LiveId = leader
			if IsFilter(&info, req, exp_strategy.Filters) {
				logs.Debug(" uid:", req.Uid, "  ", info.Uid, " LiveId:", info.LiveId)
				filter_num += 1
				continue
			}
		}
		_, isIn = hasInMap[leader]
		if isIn {
			filter_num += 1
			continue
		}
		info.Uid = uidAndLeader.Uid
		info.HasWatermark = uidAndLeader.HasWatermark
		info.FeedType = uidAndLeader.FeedType
		info.AtlasHasFilter = uidAndLeader.AtlasHasFilter
		info.Pos = uidAndLeader.Pos
		hasInMap[leader] = true
		tempList = append(tempList, info)
	}
	req.Livelist = tempList
	//logs.Debug("len:", len(req.Livelist))
}
func ConcurenceGet(req *data_type.Request, ori_url string, start, end int, routineChan chan *map[string]data_type.PhoenixResp) {
	var url bytes.Buffer
	url.WriteString(ori_url)
	for idx, info := range req.Livelist[start:end] {
		if idx != 0 {
			url.WriteString(",")
		}
		url.WriteString(info.LiveId)
	}
	result, err := discover.GetResult("phoenix", req, url.String(), 50)
	if err != nil {
		logs.Error("err:", err, " url:", url.String())
		routineChan <- nil
		return
	}
	var rMap map[string]data_type.PhoenixResp
	err = json.Unmarshal(result, &rMap)
	//logs.Debug("ConcurenceGet :", start, " end:", end, "  :", len(rMap), rMap)
	if err != nil {
		logs.Error("err:", err, "result:", string(result))
		routineChan <- nil
		return
	}
	//	tempList := make([]string, len(rMap))
	//	exp_strategy := load_confs.ExpMap[req.Mylogid]
	//	for _, info := range req.Livelist[start:end] {
	//		//是否被过滤掉
	//		leader, isIn := rMap[info.LiveId]
	//		if !isIn {
	//			continue
	//		}
	//		//是否属于需要过滤的主播
	//		info.LiveId = leader
	//		if Is_filter(info, req, exp_strategy.Filters) {
	//			continue
	//		}
	//		_, isIn = (*hasInMap)[leader]
	//		if isIn {
	//			continue
	//		}
	//		(*hasInMap)[leader] = true
	//		tempList = append(tempList, leader)
	//	}
	routineChan <- &rMap
}
