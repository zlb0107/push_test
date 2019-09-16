package feature_prepare

import (
	logs "github.com/cihub/seelog"
	//"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	//"go_common_lib/load_confs"
)

type Prepare interface {
	GetData(*data_type.Request, chan FeatureWrapperStruct) int
}

var PrepareMap map[string]Prepare

//推荐架构可配置
var PrepareCfg_map map[string]func(moduleName string, params map[string]interface{}) Prepare

func init() {
	PrepareMap = make(map[string]Prepare)

	PrepareCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) Prepare)
}

type FeatureWrapperStruct struct {
	FeatureLen      int
	NewFeatures     []NewFeatureStruct
	SnapshotVersion string
	Name            string
}
type NewFeatureStruct struct {
	//对应一个特征
	GetType     string
	Features    []RedisFeatureStruct
	FeatureType string
	FeatureLen  int
}
type RedisFeatureStruct struct {
	Dim            string
	RedisKeyPrefix string
	Features       []string
}

const (
	FEATURESNAPSHOT  string = "FeatureSnapshot"
	NORMALFEATURE    string = "NormalFeature"
	CTRNORMALFEATURE string = "CtrNormalFeature"
	CVRNORMALFEATURE string = "CvrNormalFeature"
	ONLINEFEATURE    string = "OnlineFeature"
)

func (this NewFeatureStruct) Show() {
	logs.Error("gettype:", this.GetType, " featureType:", this.FeatureType)
	for idx, rf := range this.Features {
		logs.Error("idx:", idx, "rf.Dim:", rf.Dim)
		for i, f := range rf.Features {
			logs.Error("i:", i, " f:", f)
		}
	}
}
func (this FeatureWrapperStruct) Show() {
	logs.Error("len:", this.FeatureLen, " name:", this.Name)
	logs.Error("feature_len:", len(this.NewFeatures))
	for idx, nf := range this.NewFeatures {
		logs.Error("idx:", idx, " nf:")
		logs.Error("type:", nf.GetType)
		for i, rf := range nf.Features {
			logs.Error("i:", i, " rf.Dim:", rf.Dim)
			for j, f := range rf.Features {
				logs.Error("j:", j, " f:", f)
			}
		}
	}
}
