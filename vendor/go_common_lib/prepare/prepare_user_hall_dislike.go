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

type HallDislikePrepare struct {
}

func init() {
	var rp HallDislikePrepare
	Prepare_map["HallDislikePrepare"] = rp
	logs.Warn("in raw_prepare init")
}

func (rp HallDislikePrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "HallDislikePrepare" }()
	is_hit := false
	defer common.TimerV2("HallDislikePrepare", &(request.Timer_log), time.Now(), &is_hit)
	var r Response

	ch := make(chan Response, 1)
	go func(c chan Response) {
		var r Response
		defer func() { c <- r }()

		//http://10.111.69.103:7101/api/live/recommend_dislike?type=1&uid=128758915&channel_id=1
		url := "http://10.111.149.154:8080/api/live/recommend_dislike?type=1&channel_id=0&uid=" + request.Uid
		//result, err := http_client_pool.Get_url_cache(url, "new_"+request.Uid, 3600, &is_hit)
		result, err := http_client_pool.Get_n_url(url, 10)
		if err != nil {
			logs.Error("get url:", err)
			return
		}

		if err := json.Unmarshal(result, &(r)); err != nil {
			logs.Error("Unmarshal: ", err.Error())
			return
		}
	}(ch)

	select {
	case r = <-ch:
	case <-time.After(10 * time.Millisecond):
		logs.Error("prepare::HallDislikePrepare timeout")
		return 0
	}

	if request.Dislike_hall_list == nil {
		request.Dislike_hall_list = make(map[string]bool)
	}

	for _, uid := range r.Data.Users {
		uid_str := strconv.FormatInt(uid, 10)
		request.Dislike_hall_list[uid_str] = true
	}

	return 0
}
