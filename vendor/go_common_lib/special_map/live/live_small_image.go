package live_special_map

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
	"go_common_lib/special_map"
	"math/rand"
	"time"
)

type SmallHeadImageLive struct {
	special_map.SpecialUids
}

type BlackResponse struct {
	DMErr    int      `json:"dm_error"`
	ErrorMsg string   `json:"error_msg"`
	Users    []string `json:"users"`
}

var LiveSmallImageController SmallHeadImageLive

func init() {
	go func() {
		for {
			var ips = []string{"10.111.72.122", "10.111.70.22"}
			idx := rand.Intn(len(ips))
			url := "http://" + ips[idx] + ":3125/recommond/black/ids"
			LiveSmallImageController.UpdateSimpleHttp(url, updateNew)
			time.Sleep(60 * time.Second)
		}
	}()
}

func updateNew(blackIds *map[string]bool, url string) {
	resp, err := http_client_pool.Get_url(url)
	if err != nil {
		logs.Error("error:", err)
		return
	}
	var result BlackResponse
	body_str := string(resp)
	if err := json.Unmarshal([]byte(body_str), &(result)); err != nil {
		logs.Error("Unmarshal head image result fail: ", err.Error())
		return
	}
	for _, uid := range result.Users {
		(*blackIds)[uid] = true
	}
}
