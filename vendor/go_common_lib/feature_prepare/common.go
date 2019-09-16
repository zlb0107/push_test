package feature_prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/load_models"
)

func getFeatureWrapper(conf load_models.FeatureInfo, syncChan chan int, featureArray []NewFeatureStruct, idx int, req *data_type.Request) {
	defer func() { syncChan <- idx }()
	feature := getFeature(conf, req)
	featureArray[idx] = feature
}
func getRedisKey(redisConf load_models.RedisInfo, req *data_type.Request) []interface{} {
	keys := make([]interface{}, 0)
	switch redisConf.Dim {
	case "user":
		{
			keys = append(keys, redisConf.RedisKey+req.Uid)
		}
	case "live":
		{
			for _, info := range req.Livelist {
				keys = append(keys, redisConf.RedisKey+info.Uid)
			}
		}
	case "feed":
		{
			for _, info := range req.Livelist {
				keys = append(keys, redisConf.RedisKey+info.LiveId)
			}
		}
	case "user_live":
		{
			for _, info := range req.Livelist {
				keys = append(keys, redisConf.RedisKey+req.Uid+"_"+info.Uid)
			}
		}
	case "model":
		{
			keys = append(keys, redisConf.RedisKey)
		}
	default:
		{
			logs.Error("should not be here:")
		}
	}
	return keys
}
func getRedisValue(conf load_models.RedisInfo, req *data_type.Request, waitChan chan RedisFeatureStruct) {
	rc := load_models.GetRedisClient(&conf)
	defer rc.Close()
	var redisFeature RedisFeatureStruct
	redisFeature.Dim = conf.Dim
	defer func() { waitChan <- redisFeature }()
	keys := getRedisKey(conf, req)
	if len(keys) == 0 {
		return
	}
	//	for _, key := range keys {
	//		logs.Error("key:", string(key.(string)))
	//	}

	redisFeature.RedisKeyPrefix = conf.RedisKey
	values, err := redis.Strings(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("err:", err, " len(key):", len(keys), ", redis_key:", conf.RedisKey)
		return
	}
	//	for idx, value := range values {
	//		logs.Error("key:", string(keys[idx].(string)), "value:", value)
	//	}

	//	logs.Error(len(values), " values ", len(keys), " ", len(req.Livelist))
	redisFeature.Features = values
}
func getFeature(conf load_models.FeatureInfo, req *data_type.Request) NewFeatureStruct {
	var featureInfo NewFeatureStruct
	featureInfo.GetType = conf.GetType
	featureInfo.FeatureType = conf.FeatureType
	featureInfo.FeatureLen = conf.FeatureLen
	waitNum := len(conf.Redis)
	if waitNum == 0 {
		logs.Error("redis len = 0")
		return featureInfo
	}
	featureInfo.Features = make([]RedisFeatureStruct, waitNum)
	waitChanArray := make([]chan RedisFeatureStruct, waitNum)
	for idx, single_redis := range conf.Redis {
		waitChanArray[idx] = make(chan RedisFeatureStruct, 1)
		go getRedisValue(single_redis, req, waitChanArray[idx])
	}
	for i := 0; i < waitNum; i++ {
		redisFeature := <-waitChanArray[i]
		featureInfo.Features[i] = redisFeature
	}
	return featureInfo
}
