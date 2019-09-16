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

type NewRoutineSup struct {
}

func init() {
	var rp NewRoutineSup
	NewSupMap["NewRoutineSup"] = rp
	logs.Warn("in NewRoutineSup init")
}
func (rp NewRoutineSup) Get_list(request *data_type.Request, c chan []data_type.LiveInfo, sup_info *load_confs.SupInfo) int {
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
	//redis-cli -h r-2zee626c3b9dd6e4662.redis.rds.aliyuncs.com -a vOByljzlvh26 get feed_seed_uid_167513888
	//value: 153137939300001723:click;153131166200001600:click;
	//取种子主播，无法并行
	key := sup_info.Key_prefix + request.Uid + sup_info.Key_postfix
	str, err := redis.String(rc.Do("get", key))
	if err != nil {
		logs.Error("err:", err, " key:", key)
		return -1
	}
	if len(str) == 0 {
		return 0
	}
	//该二维数组存放live_info，以便之后的深度遍历和层次遍历
	max_dep_array := 0
	//uid:source;uid:source;uid:source
	uids := strings.Split(str, ";")
	trigger_len := len(sup_info.Again)
	all_live_info := make([][]data_type.LiveInfo, 0, trigger_len)
	ch_trigger := make(chan []data_type.LiveInfo, trigger_len)
	//并发取redis
	for idx, _ := range sup_info.Again {
		go fill_one_d_info(uids, &(sup_info.Again[idx]), request, ch_trigger)
	}
	for i := 0; i < trigger_len; i++ {
		select {
		case list := <-ch_trigger:
			{
				if len(list) > max_dep_array {
					max_dep_array = len(list)
				}
				all_live_info = append(all_live_info, list)
			}
		case <-time.After(time.Millisecond * 20):
			{
				logs.Error("i:", i, " timeout")
				break
			}
		}
	}

	count := 2 * request.Count
	if sup_info.Limit_num < count {
		count = sup_info.Limit_num
	}

	list = make([]data_type.LiveInfo, 0, count)
	//进行深度或层次遍历
	num := 0
	if sup_info.Traversal == "DFS" {
		for _, infos := range all_live_info {
			for _, info := range infos {
				if num >= count {
					return 0
				}
				list = append(list, info)
				recall_num++
				num++
			}
		}
	} else {
		for j := 0; j < max_dep_array; j++ {
			for _, infos := range all_live_info {
				if num >= count {
					return 0
				}
				if j >= len(infos) {
					continue
				}
				list = append(list, infos[j])
				recall_num++
				num++
			}
		}
	}
	return 0
}
func fill_one_d_info(uids []string, sub_exp *load_confs.SupInfo, request *data_type.Request, chan_trigger chan []data_type.LiveInfo) {
	var list []data_type.LiveInfo
	defer func() { chan_trigger <- list }()
	defer common.Timer(sub_exp.Key_prefix, &(request.Timer_log), time.Now())
	var keys []interface{}
	for _, uid := range uids {
		real_uid_info := strings.Split(uid, ":")
		if len(real_uid_info) != 2 {
			logs.Error("real_uid_info is not right:", uid)
			continue
		}
		sim_key := sub_exp.Key_prefix + real_uid_info[0] + sub_exp.Key_postfix
		keys = append(keys, sim_key)
	}
	if sub_exp.Redis_client == nil {
		temp := redis_pool.RedisInfo{&(sub_exp.Redis_client), sub_exp.Redis, sub_exp.Auth, 50, 1000, 180, 10000000, 10000000}
		redis_pool.RedisInit(temp)
	}

	rc_again := sub_exp.Redis_client.Get()
	defer rc_again.Close()
	strs, err := redis.Values(rc_again.Do("mget", keys...))
	if err != nil {
		logs.Error("get redis failed:", err)
		return
	}
	exp_strategy, ok := request.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return
	}
	recall_num := 0
	noLiveNum := 0
	for i, str := range strs {
		if str == nil {
			continue
		}
		live_infos := strings.Split(string(str.([]byte)), ";")
		for _, live_info := range live_infos {
			infos := strings.Split(live_info, ":")
			if len(infos) != 2 {
				logs.Error("size is not 2:", len(infos), " key:", sub_exp.Key_prefix)
				continue
			}
			if len(list) >= sub_exp.Limit_num {
				break
			}
			recall_num++
			var info data_type.LiveInfo
			info.Uid = infos[0]
			if request.Class == "feed" {
				info.Uid = infos[0]
				info.LiveId = infos[0]
			} else {
				if info.LiveId = living.Living_handler.Get_liveid(info.Uid); info.LiveId == "" {
					noLiveNum++
					continue
				}
			}

			info.Trigger = string(keys[i].(string)) + "," + infos[1]
			score, _ := strconv.ParseFloat(infos[1], 64)
			info.Token = sub_exp.Token
			if score < sub_exp.Limit_score {
				continue
			}
			if filter.Is_filter(info, request, exp_strategy.Filters) {
				continue
			}
			info.TriggerScores = append(info.TriggerScores,
				data_type.TriggerScore{sub_exp.Key_prefix, infos[1]})
			info.Token = sub_exp.Token

			list = append(list, info)
			recall_num++
		}
	}
	GetTriggerRecallInfo(list, request.Uid, recall_num, noLiveNum, sub_exp.Key_prefix, request.EventTime)
}
