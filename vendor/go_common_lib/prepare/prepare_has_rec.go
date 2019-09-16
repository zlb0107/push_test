package prepare

import (
	"fmt"
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strings"
	"time"
)

func init() {
	Prepare_map["PrepareHasRec"] = PrepareHasRec{}
}

type PrepareHasRec struct {
}

func Gen_redis_key_has_rec(request *data_type.Request) string {
	//assert(request != nil)
	if request == nil {
		logs.Error("request is nil")
		return ""
	} else if request.Rec_tab == "" {
		//		logs.Error("request tab is empty")
		return ""
	} else if request.Uid == "" {
		//logs.Error("request uid is empty")
		return ""
	} else if request.Session_id == "" {
		//		logs.Error("request sessionId is empty")
		return ""
	} else {
		return fmt.Sprintf("%s_%s_%s", request.Rec_tab,
			request.Uid, request.Session_id)
	}
}

//fill request.Has_rec_list
func fetch_has_rec_list(key string, request *data_type.Request) int {
	rc := Has_rec_redis.Get()
	defer rc.Close()

	strs, err := redis.String(rc.Do("get", key))
	if err != nil {
		//		logs.Warn("get redis failed:", err)
		return -1
	}
	has_rec_list := make(map[string]bool)
	terms := strings.Split(strs, ";")
	for _, term := range terms {
		if term != "" {
			has_rec_list[term] = true
		}
	}
	request.Has_rec_list = has_rec_list
	return len(request.Has_rec_list)
}

func (this PrepareHasRec) Get_data(request *data_type.Request,
	ch chan string) int {
	defer func() { ch <- "PrepareHasRec" }()
	defer common.Timer("PrepareHasRec", &(request.Timer_log), time.Now())
	key := Gen_redis_key_has_rec(request)
	if key == "" {
		return -1
	}
	return fetch_has_rec_list(key, request)
}
