package load_confs

import (
	"io/ioutil"
	"strings"

	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"

	"go_common_lib/config"
	"go_common_lib/data_type"
	"go_common_lib/get_recalls"
	"go_common_lib/go-json"
)

const ExpCtxKey = "ExpKey"

type SupInfo struct {
	Id           *int
	Redis        string
	Auth         string
	Key_prefix   string
	Key_postfix  string
	Priority     int
	Traversal    string
	Limit_num    int
	Limit_score  float64
	Again        []SupInfo
	Name         string
	Token        string
	Redis_client *redis.Pool
}
type ExpStrategy struct {
	Logid               string    `json:"logid,omitempty"`
	NewSups             []SupInfo `json:"NewSups,omitempty"`
	Supplements         []string  `json:"Supplements,omitempty"`
	Merger              string    `json:"merger,omitempty"`
	Prepare             []string  `prepare strategys`
	Common_prepare      []string  `prepare strategys`
	Weight              []string  `json:"Weight,omitempty"`
	Scatter             []string  `json:"Scatter,omitempty"`
	Change              []string  `json:"Change,omitempty"`
	Interpose           []string  `json:"Interpose,omitempty"`
	Filters             []string  `json:"Filters,omitempty"`
	AttrFetcher         []string  `json:"AttrFetcher,omitempty"`
	BatchFilters        []string  `json:"BatchFilters,omitempty"`
	Sup_probility       []float64 `json:"Sup_probility,omitempty"`
	NewSupProbility     []float64 `json:"NewSupProbility,omitempty"`
	Append_probility    []float64 `json:"Append_probility,omitempty"`
	Supplements_file    string    `json:"Supplements_file,omitempty"`
	ModelKey            string    `json:"model_key,omitempty"`
	Sort                []string  `json:"Sort,omitempty"`
	Appends             []string  `json:"Appends,omitempty"`
	Models              []string  `json:"Models,omitempty"`
	ModelMerge          []string  `json:"ModelMerge,omitempty"`
	FeaturePrepare      []string  `json:"FeaturePrepare,omitempty"`
	FeatureWeight       []string  `json:"FeatureWeight,omitempty"`
	CvrModelKey         string    `json:"CvrModelKey,omitempty"`
	CtrModelKey         string    `json:"CtrModelKey,omitempty"`
	OnlineModelKey      string    `json:"OnlineModelKey,omitempty"`
	Token               string    `json:"Token,omitempty"`
	CtrXgbLeafModelKey  string    `json:"CtrXgbLeafModelKey,omitempty"`
	CvrXgbLeafModelKey  string    `json:"CvrXgbLeafModelKey,omitempty"`
	CtrDnnModelKey      string    `json:"CtrDnnModelKey,omitempty"`
	CvrDnnModelKey      string    `json:"CvrDnnModelKey,omitempty"`
	SecondModels        []string  `json:"SecondModels,omitempty"`
	SecondFeatureWeight []string  `json:"SecondFeatureWeight,omitempty"`
	NewSupMerge         string    `json:"new_sup_merge,omitempty"`

	// ModelKeyMap 配置模型的redis key，e.g. "model_key_map":{"online":"model_online_nv5_0"}}
	ModelKeyMap map[string]string `json:"model_key_map,omitempty"`

	//配置化参数的配置项
	Supplements_conf         []map[string]map[string]interface{} `json:"Supplements_conf,omitempty"`
	Prepare_conf             []map[string]map[string]interface{} `json:"Prepare_conf,omitempty"`
	Common_prepare_conf      []map[string]map[string]interface{} `json:"Common_prepare_conf,omitempty"`
	Weight_prepare_conf      []map[string]map[string]interface{} `json:"Weight_prepare_conf,omitempty"`
	Scatter_conf             []map[string]map[string]interface{} `json:"Scatter_conf,omitempty"`
	Change_conf              []map[string]map[string]interface{} `json:"Change_conf,omitempty"`
	Interpose_conf           []map[string]map[string]interface{} `json:"Interpose_conf,omitempty"`
	Filters_conf             []map[string]map[string]interface{} `json:"Filters_conf,omitempty"`
	AttrFetcher_conf         []map[string]map[string]interface{} `json:"AttrFetcher_conf,omitempty"`
	BatchFilters_conf        []map[string]map[string]interface{} `json:"BatchFilters_conf,omitempty"`
	Sort_conf                []map[string]map[string]interface{} `json:"Sort_conf,omitempty"`
	Appends_conf             []map[string]map[string]interface{} `json:"Appends_conf,omitempty"`
	Models_conf              []map[string]map[string]interface{} `json:"Models_conf,omitempty"`
	ModelMerge_conf          []map[string]map[string]interface{} `json:"ModelMerge_conf,omitempty"`
	FeaturePrepare_conf      []map[string]map[string]interface{} `json:"FeaturePrepare_conf,omitempty"`
	FeatureWeight_conf       []map[string]map[string]interface{} `json:"FeatureWeight_conf,omitempty"`
	SecondModels_conf        []map[string]map[string]interface{} `json:"SecondModels_conf,omitempty"`
	SecondFeatureWeight_conf []map[string]map[string]interface{} `json:"SecondFeatureWeight_conf,omitempty"`
	NewSup_conf              []map[string]map[string]interface{} `json:"NewSup_conf,omitempty"`

	SupplementsMerge_conf SupplementsMergeConf          `json:"SupplementsMerge_conf,omitempty"`
	Sup_probility_conf    map[string]map[string]float64 `json:"Sup_probility_conf,omitempty"`
}
type SupplementsMergeConf struct {
	FreeFlow     float64
	FlowAllocSup map[string][]string
}

type ExpConf struct {
	Logids []ExpStrategy
	Expids [][]ExpStrategy
}
type RecConf struct {
	RecTab  string `json:"rec_tab"`
	RecFile string `json:"rec_file"`
	Class   string `json:"class"`
}
type RecExpInfo struct {
	ExpMap      map[string]ExpStrategy
	NewExpArray []map[string]ExpStrategy
}

var ExpMap map[string]ExpStrategy
var NewExpArray []map[string]ExpStrategy
var UseDoubleConf bool
var ExpType string
var RecExpMap map[string]RecExpInfo
var RecClassMap map[string]string

func init() {
	//ExpType目前取值定义三种,"","userType","rec_tab"
	ExpType = config.AppConfig.String("exp::exp_type")
	RecClassMap = make(map[string]string)
	log_conf_path := config.AppConfig.String("exp::file")
	newfile := config.AppConfig.String("exp::newfile")
	if ExpType != "" {
		RecExpMap, _ = ReadDoubleConf(newfile)
		logs.Error("newfile:", newfile)
	} else {
		ExpMap, NewExpArray, _ = ReadConf(log_conf_path)
		logs.Error("expfile:", log_conf_path)
	}
}

func fillNewSups(newSups *[]SupInfo, inAgain bool) {
	for idx, newSup := range *newSups {
		if newSup.Id != nil {
			src := get_recalls.GetRecallSrc(*newSup.Id)
			if src == nil {
				logs.Error("can't found recall source, id:", *newSup.Id)
				continue
			}
			switch src.Mode {
			case 0:
				(*newSups)[idx].Name = "NoUidGetRedisSup"
			case 1:
				(*newSups)[idx].Name = "GetRedisSup"
			case 2:
				if !inAgain {
					logs.Error("mode=2 must in again array.")
				}
				(*newSups)[idx].Name = ""
			case 3:
				(*newSups)[idx].Name = src.PluginName
			}
			if src.Mode != 3 {
				(*newSups)[idx].Redis = src.RedisHost
				(*newSups)[idx].Auth = src.RedisAuth
				(*newSups)[idx].Key_prefix = src.Key
			}
			(*newSups)[idx].Token = src.Token
		}
		if len(newSup.Again) > 0 {
			fillNewSups(&(*newSups)[idx].Again, true)
		}
	}
}

func ReadConf(filename string) (map[string]ExpStrategy, []map[string]ExpStrategy, error) {
	logs.Error("filename:", filename)
	var Ec ExpConf
	expMap := make(map[string]ExpStrategy)
	newExpArray := make([]map[string]ExpStrategy, 0)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logs.Error("ReadFile: ", err.Error())
		return expMap, newExpArray, err
	}
	if err := json.Unmarshal(bytes, &Ec); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return expMap, newExpArray, err
	}
	for _, value := range Ec.Logids {
		//if len(value.Sup_probility) != len(value.Supplements) {
		if len(value.Sup_probility) != (len(value.Supplements) + len(value.Supplements_conf)) {
			panic("len is not equal:" + value.Logid)
		}
		file := value.Supplements_file
		if len(file) != 0 {
			logs.Error("file:", file)
			bytes, err := ioutil.ReadFile("./conf/" + file)
			if err != nil {
				logs.Error("ReadFile: ", err.Error())
				continue
			}
			if err := json.Unmarshal(bytes, &value.NewSups); err != nil {
				logs.Error("Unmarshal: ", err.Error())
				continue
			}
		}

		// 填充配置
		fillNewSups(&value.NewSups, false)
		expMap[value.Logid] = value
	}
	for _, ec := range Ec.Expids {
		exp_map := make(map[string]ExpStrategy)
		for _, value := range ec {
			exp_map[value.Logid] = value
		}
		newExpArray = append(newExpArray, exp_map)
	}
	logs.Error("conf:", expMap)
	logs.Error("conf:", Ec)
	return expMap, newExpArray, nil
}
func ReadDoubleConf(filename string) (map[string]RecExpInfo, error) {
	var Rc []RecConf
	recMap := make(map[string]RecExpInfo)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logs.Error("ReadFile: ", err.Error())
		return recMap, err
	}
	if err := json.Unmarshal(bytes, &Rc); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return recMap, err
	}
	for _, info := range Rc {
		RecClassMap[info.RecTab] = info.Class
		expMap, newExpArray, err := ReadConf(info.RecFile)
		if err != nil {
			logs.Error("err:", err)
			panic(info.RecFile)
		}
		recExp := RecExpInfo{ExpMap: expMap, NewExpArray: newExpArray}
		recMap[info.RecTab] = recExp
	}
	logs.Error("conf:", Rc)
	return recMap, nil
}

// GetExpStrategy 获取配置，若使用单文件配置，直接返回logid对应的配置；否则，判断是否滑屏请求，滑屏请求给rec_tab后追加"_cut"
func GetExpStrategy(req *data_type.Request) ExpStrategy {
	logid := req.Mylogid
	if UseDoubleConf == false {
		return ExpMap[logid]
	} else {
		if req.IsCut {
			var recTab string
			recTab = req.Rec_tab + "_cut"
			if _, isIn := RecExpMap[recTab]; isIn {
				return RecExpMap[recTab].ExpMap[logid]
			} else if _, isIn := RecExpMap["default_cut"]; isIn {
				return RecExpMap["default_cut"].ExpMap[logid]
			}
			logs.Error("can't find cut expStratrgy:", req.Rec_tab)
		}

		if _, isIn := RecExpMap[req.Rec_tab]; isIn {
			return RecExpMap[req.Rec_tab].ExpMap[logid]
		}
		return RecExpMap["default"].ExpMap[logid]
	}
}

//单文件配置根据singleLogid直接返回，否则根据ExpType来判断使用哪个参数.ExpType目前取值定义三种,"","userType","rec_tab"
//空字符串代表使用单文件配置,userType根据用户类型使用多文件配置,rec_tab根据request.Rec_tab使用多文件配置
func GetExpConfig(logid, userType, recTab string) ExpStrategy {
	logids := strings.Split(logid, ",")
	singleLogid := "default"
	if ExpType == "userType" {
		if _, isIn := RecExpMap[userType]; isIn {
			expMap := RecExpMap[userType].ExpMap
			singleLogid = findSingleLogid(logids, &expMap)
			return RecExpMap[userType].ExpMap[singleLogid]
		}
		return RecExpMap["default"].ExpMap[singleLogid]
	} else if ExpType == "rec_tab" {
		if _, isIn := RecExpMap[recTab]; isIn {
			exp := RecExpMap[recTab].ExpMap
			singleLogid = findSingleLogid(logids, &exp)
			return RecExpMap[recTab].ExpMap[singleLogid]
		}
		return RecExpMap["default"].ExpMap[singleLogid]
	}
	singleLogid = findSingleLogid(logids, &ExpMap)
	return ExpMap[singleLogid]
}

func findSingleLogid(logids []string, expMap *map[string]ExpStrategy) string {
	singleLogid := "default"
	for _, v := range logids {
		if v == "263" {
			singleLogid = v
			break
		}
		_, ret := (*expMap)[v]
		if !ret {
			continue
		}
		singleLogid = v
	}
	return singleLogid
}
