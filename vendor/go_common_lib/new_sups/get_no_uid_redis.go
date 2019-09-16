package new_sups

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/filters"
	"go_common_lib/living"
	"go_common_lib/load_confs"
	"go_common_lib/my_redis_pool"
	"go_common_lib/mytime"
	"strconv"
	"strings"
	"time"
)

type NoUidGetRedisSup struct {
}

func init() {
	var rp NoUidGetRedisSup
	NewSupMap["NoUidGetRedisSup"] = rp
	logs.Warn("in NoUidGetRedisSup init")
}
func (rp NoUidGetRedisSup) Get_list(request *data_type.Request, c chan []data_type.LiveInfo, sup_info *load_confs.SupInfo) int {
	var list []data_type.LiveInfo
	defer func() { c <- list }()
	defer common.Timer(sup_info.Key_prefix, &(request.Timer_log), time.Now())
	recall_num := 0
	defer func() {
		new_log := *(request.Num_log) + " " + sup_info.Key_prefix + "num:" + strconv.Itoa(recall_num)
		request.Num_log = &new_log
		//request.Num_log += " " + sup_info.Key_prefix + "num:" + strconv.Itoa(recall_num)
	}()
	if sup_info.Redis_client == nil {
		temp := redis_pool.RedisInfo{&(sup_info.Redis_client), sup_info.Redis, sup_info.Auth, 50, 1000, 180, 10000000, 10000000}
		redis_pool.RedisInit(temp)
	}
	rc := sup_info.Redis_client.Get()
	defer rc.Close()
	key := sup_info.Key_prefix
	str, err := redis.String(rc.Do("get", key))
	if err != nil {
		return 0
	}
	live_infos := strings.Split(str, ";")
	list = make([]data_type.LiveInfo, 0, len(live_infos))
	exp_strategy, ok := request.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return 0
	}
	noLiveNum := 0
	for _, live_info := range live_infos {
		infos := strings.Split(live_info, ":")
		if len(infos) != 2 {
			logs.Error("size is not 2:", len(infos), " key:", sup_info.Key_prefix)
			continue
		}
		if len(list) >= sup_info.Limit_num || len(list) >= request.Count {
			break
		}
		recall_num++
		var info data_type.LiveInfo
		info.Uid = infos[0]
		if request.Class == "feed" {
			info.LiveId = infos[0]
		} else {
			if info.LiveId = living.Living_handler.Get_liveid(info.Uid); info.LiveId == "" {
				noLiveNum++
				continue
			}
		}
		info.Trigger = key + "," + infos[1]
		info.Score, _ = strconv.ParseFloat(infos[1], 64)
		info.Token = sup_info.Token
		if info.Score < sup_info.Limit_score {
			continue
		}
		if filter.Is_filter(info, request, exp_strategy.Filters) {
			continue
		}
		info.TriggerScores = append(info.TriggerScores,
			data_type.TriggerScore{sup_info.Key_prefix, infos[1]})
		info.Token = sup_info.Token
		list = append(list, info)
	}
	GetTriggerRecallInfo(list, request.Uid, recall_num, noLiveNum, sup_info.Key_prefix, request.EventTime)
	return 0
}
