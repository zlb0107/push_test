package feed_special_map

import (
	//	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	//	"strings"
	"go_common_lib/special_map"
	"strconv"
	"time"
)

type LowQualityFeedids struct {
	special_map.SpecialUids
}

var LowQualityFeedidsController LowQualityFeedids

func init() {
	go func() {
		var keys []interface{}
		for i := 0; i < 10; i++ {
			keys = append(keys, "quality_bad_feed_"+strconv.Itoa(i))
		}
		for {
			LowQualityFeedidsController.UpdateKeys(keys, ";", LowQualityFeedRedis)
			time.Sleep(60 * time.Second)
		}
	}()
}
