// 个人电台6出1
package scatter

import (
	"go_common_lib/data_type"
	"go_common_lib/living"

	logs "github.com/cihub/seelog"
)

type ShuffleSingleRadio struct {
	BaseShuffle
}

func init() {
	var gs ShuffleSingleRadio
	Scatter_map["ShuffleSingleRadio"] = gs
	logs.Error("init ShuffleSingleRadio")
}

func (gs ShuffleSingleRadio) IsOk(req data_type.Request, idx int) (int, bool) {
	scope := 6 //多大范围内
	limit := 1 //最多几个
	count := 0

	start_idx := idx - scope + 1
	if start_idx < 0 {
		start_idx = 0
	}

	for i := start_idx; i <= idx && i < len(req.Livelist); i++ {
		liveType := living.Living_handler.Get_live_type(req.Livelist[i].Uid)
		subLiveType := living.Living_handler.GetSubLiveType(req.Livelist[i].Uid)
		if liveType == "audiolive" && subLiveType != "audiopal" {
			count += 1
		}
	}

	if count > limit {
		return scope, false
	}

	return 0, true
}

func (gs ShuffleSingleRadio) Run_shuffle(req *data_type.Request) int {
	return gs.RunShuffle(req, gs.IsOk)
}
