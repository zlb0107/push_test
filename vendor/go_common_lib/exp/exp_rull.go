package multi_exp

import (
	"context"
	"encoding/json"
	logs "github.com/cihub/seelog"
	"github.com/coreos/etcd/client"
	"io/ioutil"
	"time"
)

var Exp_handler ExpRull

type ExpRull struct {
	Levels []LevelStruct `json:"levels"`
}
type LevelStruct struct {
	Seed uint
	Exps []ExpStruct
}
type GenerateType int

const (
	UseStartEnd GenerateType = iota
	UseUids
	UseRedis
)

type ExpStruct struct {
	Expid         string
	Start         uint
	End           uint
	Priority      int
	Uids          []string
	Uids_map      map[string]bool
	Redis_addr    string
	Redis_auth    string
	Redis_key     string
	Generate_type GenerateType
}

func init() {
	raw_exp := Get_exp()
	Exp_handler, _ = Deal_raw_exp(raw_exp)
	go update()
}
func update() {
	for {
		time.Sleep(10 * time.Second)
		raw_exp := Get_exp()
		exp, ret := Deal_raw_exp(raw_exp)
		if ret != -1 {
			Exp_handler = exp
		}
		logs.Error("update exp success")
	}
}

func get_etcd_client() (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://10.111.68.176:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return client.NewKeysAPI(c), nil
}
func Deal_raw_exp(raw_exp string) (ExpRull, int) {
	var er ExpRull
	if err := json.Unmarshal([]byte(raw_exp), &er); err != nil {
		logs.Error("json open failed:", err)
		return ExpRull{}, -1
	}
	//将数组改为map，方便判断
	for i, level := range er.Levels {
		for j, exp := range level.Exps {
			er.Levels[i].Exps[j].Uids_map = make(map[string]bool)
			for _, uid := range exp.Uids {
				er.Levels[i].Exps[j].Uids_map[uid] = true
			}
		}
	}
	if 0 != Check_vaild(&er) {
		return ExpRull{}, -1
	}
	return er, 0

}
func Get_exp() string {
	defer logs.Flush()
	kapi, err := get_etcd_client()
	if err != nil {
		logs.Error("get etcd client falied:", err)
		return ""
	}
	value, err := kapi.Get(context.Background(), "/exp/generate_expid", nil)
	if err != nil {
		logs.Error("set etcd failed:", err)
		return ""
	}
	return value.Node.Value
}
func Set_exp(file string) {
	defer logs.Flush()
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logs.Error("open file failed:", err)
		return
	}
	var er ExpRull
	if err := json.Unmarshal(bytes, &er); err != nil {
		logs.Error("json open failed:", err)
		return
	}
	if 0 != Check_vaild(&er) {
		return
	}
	kapi, err := get_etcd_client()
	if err != nil {
		logs.Error("get etcd client falied:", err)
		return
	}
	_, err = kapi.Set(context.Background(), "/exp/generate_expid", string(bytes), nil)
	if err != nil {
		logs.Error("set etcd failed:", err)
	}
}
func Check_vaild(er *ExpRull) int {
	//检查各层的种子是否相同
	seed_map := make(map[uint]bool)
	for _, level := range er.Levels {
		if _, is_in := seed_map[level.Seed]; is_in {
			logs.Error("seed same:", level.Seed)
			return -1
		}
		seed_map[level.Seed] = true
	}
	//检查所有的expid是否有相同
	expid_map := make(map[string]bool)
	for _, level := range er.Levels {
		for _, exp := range level.Exps {
			if _, is_in := expid_map[exp.Expid]; is_in {
				logs.Error("expid same:", exp.Expid)
				return -1
			}
			expid_map[exp.Expid] = true
		}
	}
	//检查各层范围是否有重叠
	for _, level := range er.Levels {
		if Check_level_overlap(&level) {
			return -1
		}
	}
	return 0
}
func Check_level_overlap(level *LevelStruct) bool {
	for i, exp1 := range level.Exps {
		for j, exp2 := range level.Exps {
			if i == j {
				continue
			}
			if exp1.End < exp2.Start || exp1.Start > exp2.End {
				return false
			} else {
				logs.Error("has overlap:", exp1.Expid, ":", exp2.Expid)
				return true
			}
		}
	}
	return false
}
