package load_models

import (
	"io/ioutil"
	"time"
	"unsafe"

	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	"go_common_lib/config"
	"go_common_lib/go-json"
	"go_common_lib/my_redis_pool"
)

//key为模型的redis key，value中包括模型和特征
//小流量中，需要通过key来找对应的模型和特征
var ModelMap map[string]*FeatureConf

//最初的一个配置文件list,记录了最初的一个配置文件list
var ModelList []RedisInfo

type SingleConf struct {
	Model_list []FeatureConf
}
type ModelInfo struct {
	Redis        string
	Auth         string
	Token        string
	XgboostModel unsafe.Pointer
	Md5sum       string
	File         string
	RedisClient  *redis.Pool
	Key          string
}
type FeatureConf struct {
	NormalPb        string        `json:"normal_pb"`
	FeatureLen      int           `json:"feature_len"`
	Features        []FeatureInfo `json:"features"`
	ModelRedis      RedisInfo     `json:"model_redis"`
	ModelHandler    ModelInfo     `json:"-"`
	SnapshotVersion string        `json:"snapshot_version"`
}

type FeatureInfo struct {
	GetType string      `json:"type"`
	Redis   []RedisInfo `json:"redis"`
	//是稀疏还是正常特征
	FeatureType string `json:"feature_type"`
	FeatureLen  int    `json:"feature_len"`
}
type RedisInfo struct {
	NormalPb    string      `json:"normal_pb"`
	NewOld      string      `json:"new_old"`
	Dim         string      `json:"dim"`
	RedisKey    string      `json:"redis_key"`
	RedisAddr   string      `json:"redis_addr"`
	RedisAuth   string      `json:"redis_auth"`
	RedisClient *redis.Pool `json:"-"`
}

func init() {
	ModelMap = make(map[string]*FeatureConf)
	confFile := config.AppConfig.String("exp::confFile")
	if confFile == "" {
		logs.Error("Err: app.conf has not exp::confFile:")
		logs.Flush()
		return
	}
	ModelList = readConf(confFile)
	//以后用modellist来更新各个模型
	modelMap := update(ModelList)
	if len(modelMap) == 0 {
		logs.Error("len(modellist):", len(ModelList))
		logs.Flush()
		panic("model load failed")
	}
	ModelMap = modelMap
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			modelMap := update(ModelList)
			if len(modelMap) != 0 {
				ModelMap = modelMap
			}
		}
	}()
}

//异步更新每个模型的conf
func update(modelList []RedisInfo) map[string]*FeatureConf {
	modelMap := make(map[string]*FeatureConf)
	for _, modelConf := range modelList {
		ptrFeatureConfArray := getFeatureConfArray(modelConf)
		if ptrFeatureConfArray == nil {
			//失败
			return make(map[string]*FeatureConf)
		}
		for idx, _ := range *(ptrFeatureConfArray) {
			keyPrefix := modelConf.Dim
			var featureConf = (*ptrFeatureConfArray)[idx]
			featureConf.ModelHandler.Auth = featureConf.ModelRedis.RedisAuth
			featureConf.ModelHandler.Key = featureConf.ModelRedis.RedisKey
			featureConf.ModelHandler.Redis = featureConf.ModelRedis.RedisAddr
			featureConf.ModelHandler.File = featureConf.ModelRedis.RedisKey
			// WARN：打online快照不需要加载模型
			featureConf.NormalPb = modelConf.NormalPb
			key := keyPrefix + "_" + featureConf.ModelHandler.Key
			logs.Error("model_key:", key)
			modelMap[key] = &featureConf
		}
	}
	return modelMap
}

func getFeatureConfArray(modelConf RedisInfo) *[]FeatureConf {
	defer logs.Flush()
	for modelConf.RedisClient == nil {
		temp := redis_pool.RedisInfo_v2{&(modelConf.RedisClient), modelConf.RedisAddr, modelConf.RedisAuth, 50, 1000, 180, 1 * time.Second, 1 * time.Second, 1 * time.Second}
		redis_pool.RedisInit_v2(temp)
		logs.Error("init redis_client")
	}
	rc := modelConf.RedisClient.Get()
	defer rc.Close()
	raw_str, err := redis.String(rc.Do("get", modelConf.RedisKey))
	if err != nil {
		logs.Error("err:", err)
		return nil
	}
	var ml SingleConf
	if err := json.Unmarshal(([]byte)(raw_str), &ml); err != nil {
		logs.Error("err:", err)
		return nil
	}
	//Model_list []FeatureConf
	return &(ml.Model_list)
}

func readConf(filename string) []RedisInfo {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logs.Error("ReadFile: ", err.Error())
		panic("model file wrong")
	}
	modelList := make([]RedisInfo, 0)
	if err := json.Unmarshal(bytes, &modelList); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		panic("model file wrong:" + err.Error() + " file:" + filename)
	}
	return modelList
}
