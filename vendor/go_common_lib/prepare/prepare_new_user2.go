package prepare

import (
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/discover"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
	"go_common_lib/mytime"
)

type NewUserPrepare2 struct {
}

func init() {
	var rp NewUserPrepare2
	Prepare_map["NewUserPrepare2"] = rp
	logs.Warn("in raw_prepare init")
}

func (rp NewUserPrepare2) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "NewUserPrepare2" }()
	is_hit := false
	defer common.TimerV2("NewUserPrepare2", &(request.Timer_log), time.Now(), &is_hit)

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
	request.Is_new = nr.Is_new || nr.IsReturn30
	request.IsReturn30 = nr.IsReturn30
	return 0
}
