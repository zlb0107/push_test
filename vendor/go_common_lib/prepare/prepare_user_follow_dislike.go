package prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/http_client_pool"
	"go_common_lib/mytime"
	"strconv"
	"time"

	"go_common_lib/go-json"
)

type FollowDislikePrepare struct {
}

func init() {
	var rp FollowDislikePrepare
	Prepare_map["FollowDislikePrepare"] = rp
	logs.Warn("in raw_prepare init")
}

type Response struct {
	Data DataInfo
}
type DataInfo struct {
	Users []int64
}

func (rp FollowDislikePrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "FollowDislikePrepare" }()
	is_hit := false
	defer common.TimerV2("FollowDislikePrepare", &(request.Timer_log), time.Now(), &is_hit)
	//http://10.111.69.103:7101/api/live/recommend_dislike?type=1&uid=128758915&channel_id=1
	url := "http://10.111.149.154:8080/api/live/recommend_dislike?type=1&channel_id=6&uid=" + request.Uid
	//result, err := http_client_pool.Get_url_cache(url, "new_"+request.Uid, 3600, &is_hit)
	result, err := http_client_pool.Get_50_url(url)
	if err != nil {
		logs.Error("get url:", err)
		return -1
	}
	var r Response
	if err := json.Unmarshal(result, &(r)); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return -1
	}
	if request.Dislike_follow_list == nil {
		request.Dislike_follow_list = make(map[string]bool)
	}
	for _, uid := range r.Data.Users {
		uid_str := strconv.FormatInt(uid, 10)
		request.Dislike_follow_list[uid_str] = true
	}
	return 0
}
