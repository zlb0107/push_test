package open_falcon

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"go_common_lib/go-json"

	logs "github.com/cihub/seelog"
)

type message struct {
	Item []item `json:"item"`
}

type item struct {
	Endpoint    string  `json:"endpoint"`
	Metric      string  `json:"metric"`
	Timestamp   int64   `json:"timestamp"`
	Step        int     `json:"step"`
	Value       float64 `json:"value"`
	CounterType string  `json:"counterType"`
	Tags        string  `json:"tags"`
}

func PostToOpenFalcon(metric string, step int, value float64, tags string) error {
	apiurl := "http://127.0.0.1:1988/v1/push"
	hostname, err := os.Hostname()
	if err != nil {
		logs.Error(err)
		return err
	}

	timestamp := time.Now().Unix()

	var post message
	post.Item = append(post.Item, item{
		Endpoint:    hostname,
		Metric:      metric,
		Timestamp:   timestamp,
		Step:        step,
		Value:       value,
		CounterType: "GAUGE",
		Tags:        tags,
	})
	jsonStr, _ := json.Marshal(post.Item)
	req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
