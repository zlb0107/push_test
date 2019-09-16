package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/living"
	"sync"
)

/*
Create : 21,Aug,2019
Author : wumingqi
Modify :

Func   :
	过滤刷子视频，当在线主播数量大于某个比例时，不再过滤刷子视频
Desc   :
*/
type FilterBrushCfg struct {
	ModuleName string

	MaxRatio float64
}

var gOnce sync.Once

var gFilterBrushObserver FilterBrushCfg

//刷子视频占全量在线主播的比例
var gBrushRatio float64 = 0

func (rp FilterBrushCfg) GetModuleName() string {
	return rp.ModuleName
}

func (rp FilterBrushCfg) Update(this *living.LivingControl) error {
	brushRatio := this.GetBrushRatio()
	if brushRatio != gBrushRatio {
		logs.Warn("update gBrushRatio:", brushRatio)
	}

	gBrushRatio = brushRatio

	return nil
}

func init() {
	FilterCfg_map["FilterBrushCfg"] = NewFilterBrushCfg

	logs.Warn("in FilterBrushCfg init")
}

func NewFilterBrushCfg(moduleName string, params map[string]interface{}) Filter {
	gOnce.Do(func() {
		gFilterBrushObserver.ModuleName = "gFilterBrushObserver"
		living.Living_handler.AddObserver(gFilterBrushObserver) //添加观测者
	})

	var rp FilterBrushCfg
	rp.ModuleName = moduleName
	rp.MaxRatio = params["MaxRatio"].(float64)

	return rp
}

func (rp FilterBrushCfg) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	liveType := living.Living_handler.Get_live_type(info.Uid)
	brushNum := living.Living_handler.IsBrushVideo(info.Uid)
	if (gBrushRatio < rp.MaxRatio) && (liveType != "audiolive") && (brushNum < 0) {
		return true
	}
	return false
}
