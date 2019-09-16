// 冷启动时会先从线上同步一份配置到本地
package get_recalls

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"go_common_lib/config"
	"go_common_lib/go-json"

	logs "github.com/cihub/seelog"
	"github.com/json-iterator/go"
)

var recallSrcs sync.Map
var recallfile = config.AppConfig.String("exp::recallfile")

func init() {
	// 没有指定召回文件时，不启动更新配置操作
	if recallfile == "" {
		return
	}
	logs.Debug("init GetRecalls")
	GetRecalls()
	go func() {
		for {
			time.Sleep(time.Minute)
			GetRecalls()
		}
	}()

	go func() {
		for {
			WriteRecalls()
			time.Sleep(10 * time.Minute)
		}
	}()
}

func getRecallsFromTokenManage() (*sync.Map, error) {
	resp, err := http.Get("http://ali-c-suggest-planer05.bj:18133/recalls")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var buf Resp
	err = json.Unmarshal(body, &buf)
	if err != nil {
		return nil, err
	}

	if buf.DmError != 0 {
		return nil, errors.New("get recalls failed, err:" + buf.ErrorMsg)
	}

	var recallSrcsBuff sync.Map
	var recalls []RecallSrc
	buf.Data.ToVal(&recalls)

	for idx, src := range recalls {
		recallSrcsBuff.Store(src.Id, &recalls[idx])
	}
	return &recallSrcsBuff, nil
}

func getRecallsFromFile() (*sync.Map, error) {
	tmp, err := ioutil.ReadFile(recallfile)
	if err != nil {
		return nil, err
	}

	var recalls []RecallSrc
	err = json.Unmarshal(tmp, &recalls)
	if err != nil {
		return nil, err
	}

	var recallSrcsBuff sync.Map
	for idx, src := range recalls {
		recallSrcsBuff.Store(src.Id, &recalls[idx])
	}
	return &recallSrcsBuff, nil
}

func GetRecalls() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("GetRecalls failed, err:", err)
		}
	}()

	recallSrcsBuff, err := getRecallsFromTokenManage()
	if err != nil {
		logs.Error("GetRecalls failed, err:", err)
		recallSrcsBuff, err = getRecallsFromFile()
		if err != nil {
			panic(err)
		}
	}

	recallSrcs = *recallSrcsBuff
}

func WriteRecalls() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("WriteRecalls failed, err:", err)
		}
	}()

	var recalls []*RecallSrc
	recallSrcs.Range(func(k, v interface{}) bool {
		src := v.(*RecallSrc)
		recalls = append(recalls, src)
		return true
	})

	if len(recalls) > 0 {
		buf, err := json.Marshal(recalls)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(recallfile, buf, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func GetRecallSrc(id int) *RecallSrc {
	value, has := recallSrcs.Load(id)
	if has {
		return value.(*RecallSrc)
	}
	return nil
}

type Resp struct {
	DmError  int          `json:"dm_error"`
	ErrorMsg string       `json:"error_msg,omitempty"`
	Data     jsoniter.Any `json:"data"`
}

type RecallSrc struct {
	// Id 召回源号
	Id int `json:"id"`
	// Name 中文名称
	Name string `json:"name"`
	// RedisHost redis主机
	RedisHost string `json:"redis_host"`
	// RedisAuth reids密码
	RedisAuth string `json:"redis_auth"`
	// Key redis key
	Key string `json:"key"`
	// PluginName 自定义插件名称
	PluginName string `json:"plugin_name"`
	// ValueFormat redis value格式
	ValueFormat string `json:"value_format"`
	// Token 召回源token
	Token string `json:"token"`
	// 召回模式
	Mode int `json:"mode"`
	// Amount 召回数量
	Amount int `json:"amount"`
	// User 维护人
	User string `json:"user"`
	// Desc 说明
	Desc string `json:"desc"`
}
