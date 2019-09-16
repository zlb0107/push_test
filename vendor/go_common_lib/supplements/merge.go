package supplement

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	//	"go_common_lib/filters"
	"go_common_lib/load_confs"
	"go_common_lib/merger"
	"go_common_lib/mytime"
	"time"
)

var chance_map map[string]float64

func init() {
	chance_map = make(map[string]float64)
	chance_map[data_type.Hot_name] = 1
}
func Merge_list(request *data_type.Request, c chan data_type.ChanStruct, supplement_num int) int {
	defer common.Timer("merge_list", &(request.Timer_log), time.Now())
	if request.Livelist == nil {
		request.Livelist = make([]data_type.LiveInfo, 0)
	}
	array_map := make(map[string]*[]data_type.LiveInfo)
	//merge需要的各种map
	merge_map := make(map[string]merger.NewMergeStruct)
	//exp_strategy := load_confs.ExpMap[request.Mylogid]
	exp_strategy, ok := request.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return 0
	}
	exp_chance_map := make(map[string]float64)
	if request.Is_new {
		if len(exp_strategy.NewSupProbility) < 1 {
			exp_chance_map = chance_map
		} else {
			for idx, fp := range exp_strategy.NewSupProbility {
				exp_chance_map[exp_strategy.Supplements[idx]] = fp
			}
		}
	} else {
		if len(exp_strategy.Sup_probility) < 1 {
			exp_chance_map = chance_map
		} else {
			for idx, fp := range exp_strategy.Sup_probility {
				exp_chance_map[exp_strategy.Supplements[idx]] = fp
			}
		}

	}
	for i := 0; i < supplement_num; i++ {
		select {
		case result := <-c:
			{
				array_map[result.Name] = &result.Livelist
				var ms merger.NewMergeStruct
				ms.Length = len(result.Livelist)
				if chance, is_in := exp_chance_map[result.Name]; is_in {
					ms.Chance = chance
				}
				//没有的概率就是0，idx默认也是0
				merge_map[result.Name] = ms
			}
		case <-time.After(time.Millisecond * 80):
			{
				logs.Error("i:", i, " timeout")
				continue
			}
		}
	}
	if len(array_map) == 1 {
		//只有一个sup，不需要merge
		for _, ptr := range array_map {
			request.Livelist = *ptr
		}
		return 0
	}
	//	merger.De_duplication(array_map, merge_map, priority_list)
	liveid_map := make(map[string]data_type.LiveInfo)
	//合并
	var mergeWrapper merger.MergeWrapper
	name := mergeWrapper.Merge(merge_map)
	for name != "" {
		ms := merge_map[name]
		array := array_map[name]
		if ms.Idx >= len(*array) {
			name = mergeWrapper.Merge(merge_map)
			continue
		}
		info := (*array)[ms.Idx]
		if _, is_in := liveid_map[info.LiveId]; is_in {
			//已经重复
			ms.Idx += 1 //取下一个
			merge_map[name] = ms
			continue
		}
		if info.LiveId != "" {
			request.Livelist = append(request.Livelist, info)
		}
		liveid_map[info.LiveId] = info
		ms.Idx += 1
		merge_map[name] = ms
		//进行下一轮选择
		name = mergeWrapper.Merge(merge_map)
	}

	//初始化list
	return 0
}
