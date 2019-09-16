package feature_prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/load_confs"
	"go_common_lib/load_models"
)

func GetCommonData(req *data_type.Request, ctrOrCvr string, featureWrapper *FeatureWrapperStruct) int {
	exp_strategy, ok := req.Ctx.Value(load_confs.ExpCtxKey).(load_confs.ExpStrategy)
	if !ok {
		logs.Error("get exp_strategy fail ")
		return 0
	}
	modelKey := getModelKey(exp_strategy.ModelKeyMap, ctrOrCvr)

	featureConf, is_in := load_models.ModelMap[modelKey]
	if !is_in {
		logs.Error("ModelMap not has this redis_key:", modelKey)
		return -1
	}
	waitNum := len(featureConf.Features)
	syncChan := make(chan int, waitNum)
	req.Scope = 1 //表示新的特征，用于做快照
	featureWrapper.NewFeatures = make([]NewFeatureStruct, waitNum)
	featureWrapper.SnapshotVersion = featureConf.SnapshotVersion
	featureWrapper.FeatureLen = featureConf.FeatureLen
	for idx, conf := range featureConf.Features {
		go getFeatureWrapper(conf, syncChan, featureWrapper.NewFeatures, idx, req)
	}
	for i := 0; i < waitNum; i++ {
		<-syncChan
	}
	logs.Flush()
	return 0
}

func getModelKey(modelKeyMap map[string]string, ctrOrCvr string) string {
	if modelKey, has := modelKeyMap[ctrOrCvr]; has {
		return ctrOrCvr + "_" + modelKey
	}
	logs.Error("not has this redis_key:", modelKeyMap, "ctrOrCvr:", ctrOrCvr)
	return ""
}
