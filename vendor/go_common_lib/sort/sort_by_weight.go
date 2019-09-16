package sort

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"sort"
	"strconv"
	"time"
)

type SortByWeight struct {
}

func init() {
	var rp SortByWeight
	SortMap["SortByWeight"] = rp
	logs.Warn("in distance level sort init")
}

type SortWeightType []data_type.LiveInfo

func (c SortWeightType) Len() int {
	return len(c)
}
func (c SortWeightType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c SortWeightType) Less(i, j int) bool {
	return c[i].Score > c[j].Score
}

func (rp SortByWeight) Run_sort(request *data_type.Request) int {
	defer common.Timer("SortByWeight", &(request.Timer_log), time.Now())
	sort.Sort(SortWeightType(request.Livelist))
	for i, _ := range request.Livelist {
		request.Livelist[i].Token += strconv.Itoa(i)
	}
	return 0
}
