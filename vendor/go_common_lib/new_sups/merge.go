package new_sups

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/mytime"
	"runtime"
	"strconv"
	"time"
	//"feed_thunder/connect_pool/gen-go/living"
	"bytes"
	"encoding/json"
	"go_common_lib/data_type"
	"go_common_lib/load_confs"
	"go_common_lib/snapshot"
)

var NewSupMergeMap map[string]Merge

type Merge interface {
	Merge(request *data_type.Request, c chan []data_type.LiveInfo, supplement_num int) int
}

func init() {
	NewSupMergeMap = make(map[string]Merge)
}

type TokenMerge struct {
	Info   data_type.LiveInfo
	Tokens map[string]int
}
type BigDataInfo struct {
	EventTime  string            `json:"event_time"`
	EventTopic string            `json:"event_topic"`
	EventUid   string            `json:"event_uid"`
	EventPlace string            `json:"event_place"`
	Atom       string            `json:"atom"`
	Info       TriggerRecallInfo `json:"info"`
}
type TriggerRecallInfo struct {
	Uid        string `json:"uid"`
	Trigger    string `json:"trigger"`
	RecallNum  int    `json:"recall_num"`
	FinalNum   int    `json:"final_num"`
	NotLiveNum int    `json:"not_live_num"`
	NaT        string `json:"natime"`
	Liveuids   string `json:"liveuids"`
}

func GetTriggerRecallInfo(list []data_type.LiveInfo, uid string, recallNum int, notLiveNum int, trigger, eventTime string) {
	if len(uid) < 4 {
		return
	}
	if uid[3] != '7' {
		return
	}
	var buf bytes.Buffer
	for idx, info := range list {
		if idx != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(info.Uid)
		buf.WriteString("_")
		buf.WriteString(info.LiveId)
	}
	triggerRecallInfo := TriggerRecallInfo{Uid: uid, Trigger: trigger, RecallNum: recallNum, FinalNum: len(list), NotLiveNum: notLiveNum, Liveuids: buf.String(), NaT: eventTime}
	var bigDataInfo BigDataInfo
	bigDataInfo.EventTime = eventTime
	bigDataInfo.EventTopic = "rechall_trigger_details"
	bigDataInfo.EventUid = uid
	bigDataInfo.EventPlace = "10.111.68.136"
	bigDataInfo.Info = triggerRecallInfo
	infoStr, _ := json.Marshal(bigDataInfo)
	if runtime.NumGoroutine() > 5000 {
		logs.Error("len(Channel):", len(snapshot.ChannelTrigger), " throw out this snapshot")
		//做降级
		return
	}
	snapshot.ChannelTrigger <- string(infoStr)
}

func Merge_list(request *data_type.Request, c chan []data_type.LiveInfo, supplement_num int) int {
	exp_strategy, ok := request.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return 0
	}
	if exp_strategy.NewSupMerge != "" {
		inter, ok := NewSupMergeMap[exp_strategy.NewSupMerge]
		if ok {
			return inter.Merge(request, c, supplement_num)
		}
	}

	return Merge_list_default(request, c, supplement_num)
}

func Merge_list_default(request *data_type.Request, c chan []data_type.LiveInfo, supplement_num int) int {
	defer common.Timer("sup_merge", &(request.Timer_log), time.Now())
	recall_num := 0
	defer func() {
		new_log := *(request.Num_log) + " sup_merge_num:" + strconv.Itoa(recall_num)
		request.Num_log = &new_log
	}()

	//	var living_req living.LivingContent
	uid_map := make(map[string]TokenMerge)
	for i := 0; i < supplement_num; i++ {
		select {
		case list := <-c:
			{
				for _, info := range list {
					uid := info.LiveId
					if tm, ok := uid_map[uid]; ok {
						if info.RecReason == "plan" {
							tm.Info = info
							uid_map[uid] = tm
						} else {
							//已经出现该uid，将token加入tokens
							tm.Tokens[info.Token] = 1
						}
						//累加trigger
						tm.Info.Trigger += "^" + info.Trigger
						tm.Info.TriggerScores = append(tm.Info.TriggerScores,
							info.TriggerScores...)

						uid_map[uid] = tm

					} else {
						var tm TokenMerge
						tm.Tokens = make(map[string]int)
						tm.Tokens[info.Token] = 1
						tm.Info = info
						uid_map[uid] = tm
					}
				}
			}
		case <-time.After(time.Millisecond * 20):
			{
				logs.Error("i:", i, " timeout")
				continue
			}
		}
	}
	//初始化list
	if request.Livelist == nil {
		request.Livelist = make([]data_type.LiveInfo, 0)
	}
	for _, tm := range uid_map {
		//rec_7_5_1_0^470722321_1503676802811_24^1,2,3,4
		info := tm.Info
		idx := 0
		tstr := "^"
		if info.RecReason == "plan" {
			info.Append_token = info.Token
		} else {
			//计划的append_token是用来表明plan——id的
			for token, _ := range tm.Tokens {
				if idx == 0 {
					tstr += token
				} else {
					tstr += "," + token
				}
				idx++
			}
			info.Append_token += tstr
		}
		request.Livelist = append(request.Livelist, info)
	}
	recall_num = len(request.Livelist)
	return 0
}
