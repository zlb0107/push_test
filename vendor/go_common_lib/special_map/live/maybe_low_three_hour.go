package live_special_map

import (
	//	logs "github.com/cihub/seelog"
	//	"github.com/garyburd/redigo/redis"
	//	"strings"
	"go_common_lib/special_map"
	"time"
)

type MaybeLowThreeHour struct {
	special_map.SpecialUids
}

var MaybeLowThreeHourController MaybeLowThreeHour

//redis-cli -h r-2ze9a3012577b554616.redis.rds.aliyuncs.com -a m1GwbBzf6uvm get check_filter_hour_3
func init() {
	go func() {
		for {
			MaybeLowThreeHourController.Update("check_filter_hour_3", ";", MaybeLowRedis)
			time.Sleep(60 * time.Second)
		}
	}()
}
