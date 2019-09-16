package scatter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
)

type FeedUidShuffle struct {
	BaseShuffle
}

func init() {
	var gs FeedUidShuffle
	Scatter_map["FeedUidShuffle"] = gs
	logs.Error("init FeedUidShuffle")
}
func (gs FeedUidShuffle) IsOk(req data_type.Request, idx int) (int, bool) {
	scope := 60 //多大范围内
	limit := 1  //最多几个
	count := 0
	start_idx := idx - scope + 1
	if start_idx < 0 {
		start_idx = 0
	}
	uidsMap := make(map[string]bool)
	for i := start_idx; i <= idx && i < len(req.Livelist); i++ {
		_, isIn := uidsMap[req.Livelist[i].Uid]
		if isIn {
			count += 1
		} else {
			uidsMap[req.Livelist[i].Uid] = true
		}
	}
	if count >= limit {
		return scope, false
	}
	return 0, true
}
func (gs FeedUidShuffle) Run_shuffle(req *data_type.Request) int {
	return gs.RunShuffle(req, gs.IsOk)
}
