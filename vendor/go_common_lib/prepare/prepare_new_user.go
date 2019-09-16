package prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/discover"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
	"go_common_lib/mytime"
	"time"
)

type NewUserPrepare struct {
}

func init() {
	var rp NewUserPrepare
	Prepare_map["NewUserPrepare"] = rp
	logs.Warn("in raw_prepare init")
}

type NewResponse struct {
	Is_new     bool
	IsReturn30 bool `json:"is_return30"`
}

func (rp NewUserPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "NewUserPrepare" }()
	is_hit := false
	defer common.TimerV2("NewUserPrepare", &(request.Timer_log), time.Now(), &is_hit)
	// ip := discover.GetUrl("diting", request)
	// //url := "http://10.111.95.164:18095/diting?uid=" + request.Uid
	// url := "http://" + ip + ":18095/diting?uid=" + request.Uid
	// result, err := http_client_pool.Get_url_cache(url, "new_"+request.Uid, 3600, &is_hit)
	// //result, err := http_client_pool.Get_url(url)
	// if err != nil {
	// 	logs.Error("get url:", err)
	// 	return -1
	// }

	var result []byte
	url := ":18095/diting?uid=" + request.Uid
	cacheKey := "new_" + request.Uid

	result, is_hit = http_client_pool.GetUrlResultFromCache(cacheKey)
	if !is_hit {
		var err error
		result, err = discover.GetResult("diting", request, url, 20)
		if err != nil {
			logs.Error("get url:", err)
			return -1
		}

		http_client_pool.UpdateUrlResultCache(cacheKey, string(result), 3600)
	}

	var nr NewResponse
	if err := json.Unmarshal(result, &(nr)); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return -1
	}
	request.Is_new = nr.Is_new
	request.IsReturn30 = nr.IsReturn30
	return 0
}
