package live_special_map

import (
	//	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	//	"strings"
	"go_common_lib/special_map"
	"time"
)

type PerRecBlackUids struct {
	special_map.SpecialUids
}

var PerRecBlackController PerRecBlackUids

func init() {
	go func() {
		for {
			//495328314:1532188800;621943489:1532189280;
			PerRecBlackController.UpdateTwoLevel("rec_black_uid", ";", ":", PerRecBlackRedis)
			time.Sleep(60 * time.Second)
		}
	}()
}
