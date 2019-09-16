package controllers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go_common_lib/async_task"
	"go_common_lib/change"
	"go_common_lib/data_type"
	"go_common_lib/exp"
	"go_common_lib/feature_prepare"
	"go_common_lib/feature_weight"
	"go_common_lib/fetchattr"
	"go_common_lib/filters"
	"go_common_lib/go-json"
	"go_common_lib/interpose"
	"go_common_lib/load_confs"
	"go_common_lib/model_merge"
	"go_common_lib/models"
	"go_common_lib/mytime"
	"go_common_lib/new_sups"
	"go_common_lib/prepare"
	common_lib "go_common_lib/prepare"
	"go_common_lib/scatter"
	"go_common_lib/sort"
	"go_common_lib/supplements"

	logs "github.com/cihub/seelog"
)

type LogicController struct {
	req *http.Request
}

/*
   1:输入逻辑处理
   2:小流量控制
   3:执行组合逻辑
   	 1)召回
	 2)过滤
	 3)打散
   4:输出逻辑处理
*/
func (this *LogicController) Stack() []byte {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return buf[:n]
}

//func (this *LogicController) Get_parameter() (int, data_type.Request) {
//	start := time.Now()
//	uid := this.GetString("uid")
//	logid := this.GetString("logid")
//	if uid == "" || len(uid) > 10 {
//		logs.Error("this request has not uid:", uid)
//		return -1, data_type.Request{}
//	}
//	//  if logid == "" {
//	//      logs.Warn("this request has not logid")
//	//  }
//	var re data_type.Request
//	re.Uid = uid
//	re.Logid = logid
//	re.Session_id = this.GetString("session_id")
//	re.Gender = this.GetString("gender")
//	re.Count, _ = this.GetInt("count")
//	re.Page_idx, _ = this.GetInt("page_idx")
//	re.Rec_tab = this.GetString("rec_tab")
//	re.Latitude = this.GetString("latitude")
//	re.Longitude = this.GetString("longitude")
//	re.UserLevel, _ = this.GetInt("user_level")
//	if re.Count == 0 {
//		//TODO,接入分页后，逻辑需要判断请求类型了
//		re.Count = 200
//	}
//	re.Real_count = re.Count
//	new_log := ""
//	re.Num_log = &new_log
//	common.Timer("get_parameter", &(re.Timer_log), start)
//	return 0, re
//}

func (this *LogicController) Flow(getPara func() (int, data_type.Request), outPut func(*data_type.Request) *[]byte) *[]byte {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Controller: Panic. panic message: %#v. stack info: \n%s", r, this.Stack())
			logs.Error(msg)
			logs.Flush()
		}
	}()

	start := time.Now()
	//step 1
	//判定返回状态值, 0为成功，其余负值为失败
	ret, req := getPara()
	if ret != 0 {
		return nil
	}
	req.AttrMap = make(map[string]interface{}, 0)
	req.OtherAttrMap = make(map[string]interface{}, 0)

	if 0 != this.Deal(&req) {
		return nil
	}
	//step 4
	req_str_ptr := outPut(&req)
	dis := time.Since(start).Nanoseconds()
	msecond := dis / 1000000
	new_log := ""
	if req.Timer_log != nil {
		new_log = *(req.Timer_log)
	}
	num_log := ""
	if req.Num_log != nil {
		num_log = *(req.Num_log)
	}

	ip := this.GetIP()

	//logs.Warn("list_size:", len(req.Livelist), " total:", strconv.FormatInt(msecond, 10), new_log, " ", num_log)
	logs.Warn("uid:", req.Uid, " logid:", req.Mylogid, " rec_tab:", req.Rec_tab, " sessionID:", req.Session_id, " count:", req.Count, " list_size:", len(req.Livelist), " raw_len:", req.Raw_len, " total:", strconv.FormatInt(msecond, 10), new_log, num_log, " longitude:", req.Longitude, " latitude:", req.Latitude, " ip:", ip)
	return req_str_ptr
}
func (this *LogicController) Deal(req *data_type.Request) int {
	//step 2
	//his.find_my_logid(req)
	//exp_strategy := load_confs.ExpMap[req.Mylogid]
	exp_strategy := load_confs.GetExpConfig(req.Logid, "", req.Rec_tab)
	req.Mylogid = exp_strategy.Logid
	req.Ctx = context.WithValue(context.Background(), load_confs.ExpCtxKey, exp_strategy)
	//step 3
	//这里的同步用sync.WaitGroup更好，可惜已重度依赖，改起来收益不大
	common_prepare_len := len(exp_strategy.Common_prepare)
	chs_common_prepare := make([]chan string, common_prepare_len)
	for i, interface_name := range exp_strategy.Common_prepare {
		inter, ret := common_lib.Prepare_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		chs_common_prepare[i] = make(chan string)
		go inter.Get_data(req, chs_common_prepare[i])
	}
	for i := 0; i < common_prepare_len; i++ {
		<-chs_common_prepare[i]
	}

	newSupLen := len(exp_strategy.NewSups)
	chNewSup := make(chan []data_type.LiveInfo, newSupLen)
	for idx, interface_exp := range exp_strategy.NewSups {
		inter, ret := new_sups.NewSupMap[interface_exp.Name]
		if !ret {
			logs.Error("not has this name:", interface_exp.Name)
			continue
		}
		go inter.Get_list(req, chNewSup, &(exp_strategy.NewSups[idx]))
	}
	new_sups.Merge_list(req, chNewSup, newSupLen)

	//3.0
	supplement_len := len(exp_strategy.Supplements)
	ch_supplement := make(chan data_type.ChanStruct, supplement_len)
	for _, interface_name := range exp_strategy.Supplements {
		inter, ret := supplement.Supplement_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		go inter.Get_list(req, ch_supplement)
	}
	supplement.Merge_list(req, ch_supplement, supplement_len)
	*(req.Num_log) += " after_sup_num:" + strconv.Itoa(len(req.Livelist))

	changeLen := len(exp_strategy.Change)
	chs_change := make([]chan string, changeLen)
	for i, interface_name := range exp_strategy.Change {
		inter, ret := change.Change_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		chs_change[i] = make(chan string)
		go inter.ChangeData(req, chs_change[i])
	}
	for i := 0; i < changeLen; i++ {
		<-chs_change[i]
	}

	//批量过滤
	for _, interface_name := range exp_strategy.BatchFilters {
		inter, ret := filter.BatchFilterMap[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		inter.FilterInfos(req)
	}

	//3.1
	prepare_len := len(exp_strategy.Prepare)
	chs_prepare := make([]chan string, prepare_len)
	for i, interface_name := range exp_strategy.Prepare {
		inter, ret := prepare.Prepare_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		chs_prepare[i] = make(chan string)
		go inter.Get_data(req, chs_prepare[i])
	}
	for i := 0; i < prepare_len; i++ {
		<-chs_prepare[i]
	}
	//feature_prepare
	//解决新特征的提取问题，以及快照问题。
	//第一层并发，新特征与下版本特征快照并发
	//第二层并发，取各个特征（新特征还是下版本特征并发）
	//之后聚合到一个结构中
	featurePrepareLen := len(exp_strategy.FeaturePrepare)
	featurePrepareChan := make(chan feature_prepare.FeatureWrapperStruct, featurePrepareLen)
	featureRecords := make(map[string]feature_prepare.FeatureWrapperStruct)
	for _, interfaceName := range exp_strategy.FeaturePrepare {
		inter, ret := feature_prepare.PrepareMap[interfaceName]
		if !ret {
			logs.Error("not has this name:", interfaceName)
			continue
		}
		go inter.GetData(req, featurePrepareChan)
	}
	for i := 0; i < featurePrepareLen; i++ {
		featureRecond := <-featurePrepareChan
		featureRecords[featureRecond.Name] = featureRecond
	}
	//组合特征
	featureWeightLen := len(exp_strategy.FeatureWeight)
	featureWeightChan := make(chan string, featureWeightLen)
	for _, interface_name := range exp_strategy.FeatureWeight {
		inter, ret := feature_weight.Weight_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		go inter.Get_weight(req, &featureRecords, featureWeightChan)
	}
	for i := 0; i < featureWeightLen; i++ {
		<-featureWeightChan
	}
	modelsLen := len(exp_strategy.Models)
	modelsChan := make(chan string, modelsLen)
	for _, interfaceName := range exp_strategy.Models {
		inter, ret := models.ModelsMap[interfaceName]
		if !ret {
			logs.Error("not has this name:", interfaceName)
			continue
		}
		go inter.Predict(req, &featureRecords, modelsChan)
	}
	for i := 0; i < modelsLen; i++ {
		<-modelsChan
	}

	secondFeatureWeightLen := len(exp_strategy.SecondFeatureWeight)
	secondFeatureWeightChan := make(chan string, secondFeatureWeightLen)
	for _, interface_name := range exp_strategy.SecondFeatureWeight {
		inter, ret := feature_weight.Weight_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		go inter.Get_weight(req, &featureRecords, secondFeatureWeightChan)
	}
	for i := 0; i < secondFeatureWeightLen; i++ {
		<-secondFeatureWeightChan
	}

	secondModelsLen := len(exp_strategy.SecondModels)
	secondModelsChan := make(chan string, secondModelsLen)
	for _, interfaceName := range exp_strategy.SecondModels {
		inter, ret := models.ModelsMap[interfaceName]
		if !ret {
			logs.Error("SecondModels not has this name:", interfaceName)
			continue
		}
		go inter.Predict(req, &featureRecords, secondModelsChan)
	}
	for i := 0; i < secondModelsLen; i++ {
		<-secondModelsChan
	}

	//model
	for _, interface_name := range exp_strategy.ModelMerge {
		inter, ret := model_merge.ModelMergeMap[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		inter.Merge(req)
	}
	//sort
	for _, interface_name := range exp_strategy.Sort {
		inter, ret := sort.SortMap[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		inter.Run_sort(req)
	}

	//scatter
	for _, interface_name := range exp_strategy.Scatter {
		inter, ret := scatter.Scatter_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		inter.Run_shuffle(req)
	}

	//interpose
	for _, interface_name := range exp_strategy.Interpose {
		inter, ret := interpose.Interpose_map[interface_name]
		if !ret {
			logs.Error("not has this name:", interface_name)
			continue
		}
		inter.Run_interpose(req)
	}
	*(req.Num_log) += " before_cut_num:" + strconv.Itoa(len(req.Livelist))
	if len(req.Livelist) > req.Count {
		req.Livelist = append(req.Livelist[:req.Count])
	}
	*(req.Num_log) += " after_cut_num:" + strconv.Itoa(len(req.Livelist))
	//3.5 产生token
	this.gen_token(req)
	//3.6 attr fetch
	fetch_real_count := len(exp_strategy.AttrFetcher)
	fetch_channel := make(chan string, fetch_real_count)
	for _, fetch_name := range exp_strategy.AttrFetcher {
		inter, ret := fetchattr.Fetcher_map[fetch_name]
		if !ret {
			logs.Error("not has this fetcher:", fetch_name)
			fetch_real_count--
			continue
		}
		inter.Get_attr(req, fetch_channel)
	}
	for i := 0; i < fetch_real_count; i++ {
		<-fetch_channel
	}
	//记录本次推荐结果
	if len(req.Rec_tab) >= 4 && req.Rec_tab[:4] == "feed" {
		if !async_task.AddAsyncTask(req, async_task.ReadyQueue_feed_record_has_rec) {
			logs.Warn("add request to ReadyQueue_record_has_rec failed")
		}
	} else {

		if !async_task.AddAsyncTask(req, async_task.ReadyQueue_record_has_rec) {
			logs.Warn("add request to ReadyQueue_record_has_rec failed")
		}
	}

	return 0
}

func (this *LogicController) gen_token(req *data_type.Request) {
	for idx, info := range req.Livelist {
		req.Livelist[idx].Token += "^" + strconv.Itoa(idx) + info.Append_token
	}

}
func (this *LogicController) trans_out(req *data_type.Request) *[]byte {
	defer common.Timer("out", &(req.Timer_log), time.Now())
	//填充推荐理由
	for i, _ := range req.Livelist {
		req.Livelist[i].Other_para.Reason = req.Livelist[i].RecReason
	}
	req_str, _ := json.Marshal(req.Livelist)

	return &req_str
}
func (this *LogicController) find_my_logid(req *data_type.Request) {
	if load_confs.UseDoubleConf == false {
		this.find_single_logid(req, &load_confs.ExpMap, &load_confs.NewExpArray)
	} else {
		recTab := "default"
		if _, isIn := load_confs.RecExpMap[req.Rec_tab]; isIn {
			recTab = req.Rec_tab
		}
		recExp := load_confs.RecExpMap[recTab]
		this.find_single_logid(req, &(recExp.ExpMap), &(recExp.NewExpArray))
	}
}

func (this *LogicController) find_single_logid(req *data_type.Request,
	expMap *map[string]load_confs.ExpStrategy, newExpArray *[]map[string]load_confs.ExpStrategy) {
	logid := "default"
	logids := strings.Split(req.Logid, ",")
	for _, v := range logids {
		if v == "263" {
			logid = v
			break
		}
		_, ret := (*expMap)[v]
		if !ret {
			continue
		}
		logid = v
	}
	req.Mylogid = logid
	expids_str := multi_exp.Get_expids(req.Uid)
	expids := strings.Split(expids_str, ",")
	if req.Expids == nil {
		req.Expids = make([]string, 0)
	}
	for _, expid := range expids {
		//NewExpArray []map[string]ExpStrategy
		for _, exp_map := range *newExpArray {
			if _, is_in := exp_map[expid]; is_in {
				req.Expids = append(req.Expids, expid)
			}
		}
	}
}

func (this *LogicController) GetString(name string) string {
	return this.req.URL.Query().Get(name)
}

func (this *LogicController) GetBool(name string, def ...bool) (bool, error) {
	val := this.req.URL.Query().Get(name)
	if len(val) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseBool(val)
}

func (this *LogicController) GetInt(name string) (int, error) {
	return strconv.Atoi(this.req.URL.Query().Get(name))
}

func (this *LogicController) SetRequest(req *http.Request) {
	this.req = req
}

func (this *LogicController) GetIP() string {
	ips := this.req.Header.Get("X-Forwarded-For")
	if ips != "" {
		return strings.Split(ips, ",")[0]
	}

	if ip, _, err := net.SplitHostPort(this.req.RemoteAddr); err != nil {
		return ip
	}

	return this.req.RemoteAddr
}
