package feature_weight

import (
	"go_common_lib/data_type"
	"go_common_lib/feature_prepare"
)

type Weight interface {
	Get_weight(*data_type.Request, *map[string]feature_prepare.FeatureWrapperStruct, chan string) int
}

var Weight_map map[string]Weight

//推荐架构可配置
var WeightCfg_map map[string]func(moduleName string, params map[string]interface{}) Weight

func init() {
	Weight_map = make(map[string]Weight)

	WeightCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Weight)
}
func delete_liveid(request *data_type.Request, i int) {
	request.Livelist = append(request.Livelist[:i], request.Livelist[i+1:]...)
}

type DealInfo interface {
	DealUserInfo(map[int]UserCacheWrapper, feature_prepare.RedisFeatureStruct, int, *data_type.LiveInfo, string, int) [][]float64
	DealLiveInfo(int, feature_prepare.RedisFeatureStruct, *data_type.LiveInfo, string, int) [][]float64
}
type UserCacheWrapper struct {
	Snapshot string
	Features [][]float64
}

//func putSnapToOffline(snap *string, liveInfo *data_type.LiveInfo) {
//	if len(liveInfo.OfflineSnapshot) != 0 {
//		liveInfo.OfflineSnapshot += "^"
//	}
//	liveInfo.OfflineSnapshot += *snap
//}
//func putSnapToOnline(snap *string, liveInfo *data_type.LiveInfo) {
//	if len(liveInfo.OnlineSnapshot) != 0 {
//		liveInfo.OnlineSnapshot += "^"
//	}
//	liveInfo.OnlineSnapshot += *snap
//}
