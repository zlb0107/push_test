package discover

import (
	"strconv"

	"go_common_lib/http_client_pool"
)

var (
	TimeChan     = make(chan SnameTimeInfo, 10000)
	TotalTimeMap = make(map[string]int64)
	CountMap     = make(map[string]int)
)

func init() {
	go func() {
		DealTime()
	}()
}

func DealTime() {
	for {
		timeInfo := <-TimeChan
		key := timeInfo.sname + ":" + timeInfo.ip
		total, isIn := TotalTimeMap[key]
		if isIn {
			total += int64(timeInfo.time)
		} else {
			total = int64(timeInfo.time)
		}
		TotalTimeMap[key] = total
		count, cIsIn := CountMap[key]
		if cIsIn {
			count += 1
		} else {
			count = 1
		}
		CountMap[key] = count
		if count >= 10 {
			avrg := int(total) / count
			//通知,写入调用接口
			url := "http://10.111.95.183:18097/discover/report?server_name=" + timeInfo.sname + "&ip=" + timeInfo.ip + "&time=" + strconv.Itoa(avrg)
			http_client_pool.Get_url(url)
			TotalTimeMap[key], CountMap[key] = 0, 0
		}
	}
}
