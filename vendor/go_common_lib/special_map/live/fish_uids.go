package live_special_map

import (
	//	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	//	"strings"
	"go_common_lib/special_map"
	"time"
)

type FishUids struct {
	special_map.SpecialUids
}

var FishUidsController FishUids

func init() {
	go func() {
		for {
			FishUidsController.Update("total_fish_label_uid", ";", FishRedis)
			time.Sleep(3600 * time.Second)
		}
	}()
}
