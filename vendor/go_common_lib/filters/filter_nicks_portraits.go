package filter

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"go_common_lib/data_type"
)

type FilterNicksPortraits struct {
}

func init() {
	var rp FilterNicksPortraits
	BatchFilterMap["FilterNicksPortraits"] = rp
	logs.Warn("FilterNicksPortraits init")
}
func (rp FilterNicksPortraits) FilterInfos(request *data_type.Request) {
	var keys []interface{}
	// redis-cli -h r-2ze5ec929f2957d4.redis.rds.aliyuncs.com -p 6379 -a Ne4w1Riy3 get fai_10006135
	//value：1_1 第一位表示昵称，第二位表示头像，1表示有，0表示无
	for _, info := range request.Livelist {
		keys = append(keys, "fai_"+info.LiveId)
	}
	rc := NicksPortraitsRedis.Get()
	defer rc.Close()
	values, err := redis.Strings(rc.Do("mget", keys...))
	if err != nil {
		logs.Error("err:", err)
		return
	}
	tempList := make([]data_type.LiveInfo, 0)
	for idx, value := range values {
		if value != "" && value != "1_1" {
			logs.Error("feedid:", request.Livelist[idx].LiveId, " status:", value)
			continue
		}
		tempList = append(tempList, request.Livelist[idx])
	}
	request.Livelist = tempList
}
