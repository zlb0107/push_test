package prepare

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	//"go_common_lib/cache"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strings"
	"time"
)

type FilmHasShownPrepare struct {
}

func init() {
	var rp FilmHasShownPrepare
	Prepare_map["FilmHasShownPrepare"] = rp
	logs.Warn("in FilmHasShownPrepare init")
	//cache.Cache_controlor.Register("rec_filter_")
}
func (rp FilmHasShownPrepare) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "FilmHasShownPrepare" }()
	//	defer common.Timer("FilmHasShownPrepare", &(request.Timer_log), time.Now())
	is_hit := false
	defer common.TimerV2("FilmHasShownPrepare", &(request.Timer_log), time.Now(), &is_hit)
	rc := FilmHasShownRedis.Get()
	defer rc.Close()
	strs, err := redis.String(rc.Do("get", "rec_film_filter_"+request.Uid))
	//strs, err := cache.Get_redis_string("rec_filter_"+request.Uid, rc, &is_hit, 60)
	if err != nil {
		//logs.Error("get redis failed:", err, " uid:", request.Uid)
		//logs.Flush()
		return -1
	}
	//4836893:-0.5;88341318:-0.5;4296493:-0.5
	uid_shown_map := make(map[string]bool)
	terms := strings.Split(strs, ";")
	for _, term := range terms {
		uid_shown_map[term] = true
		//	logs.Error("has_shown:", term)
	}
	request.Film_has_shown_list = uid_shown_map

	return 0
}
