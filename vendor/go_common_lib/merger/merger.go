package merger

import (
	"math/rand"
	"time"
)

type MergeStruct struct {
	Length   int     //队列长度
	Idx      int     //队列现偏移
	Chance   float64 //选取该队列的概率
	Rand_max float64 //随机范围最大值
	Rand_min float64 //随机范围最小值
}

func Merge(Merge_map map[string]MergeStruct) string {
	//先调整cmap，概率发生变化原因，有些队列已遍历完
	need_change := false
	total := .0
	for _, ms := range Merge_map {
		if ms.Idx < ms.Length {
			total += ms.Chance
			continue
		}
		need_change = true
	}
	if need_change {
		for key, ms := range Merge_map {
			if ms.Idx < ms.Length {
				ms.Chance = ms.Chance / total
			} else {
				//已经取光的队列
				ms.Chance = 0
			}
			(Merge_map)[key] = ms
		}
	}
	//重新计算随机范围
	cumulation := 0.0
	i := 0
	for key, ms := range Merge_map {
		i += 1
		ms.Rand_min = cumulation
		cumulation += ms.Chance
		ms.Rand_max = cumulation
		if i == len(Merge_map) {
			//最后一个元素
			ms.Rand_max = 1.0
		}
		if ms.Rand_min == .0 && ms.Rand_max == 1.0 {
			//只剩它了
			if ms.Idx >= ms.Length {
				return ""
			}
			return key
		}
		(Merge_map)[key] = ms
	}
	//进行随机策略，选取队列
	//将各个候选队列划分为0-1的float候选区，随机数落在哪个区，就属于哪个队列
	rand.Seed(time.Now().UnixNano())
	randnum := rand.Float64()
	//寻找命中区间，左闭右开
	for key, ms := range Merge_map {
		if randnum >= ms.Rand_min && randnum < ms.Rand_max {
			return key
		}
	}
	return ""
}
