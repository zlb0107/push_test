package interpose

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"time"
)

type TimeSectionFeed struct {
}

func init() {
	var rp TimeSectionFeed
	Interpose_map["TimeSectionFeed"] = rp
	logs.Warn("in TimeSectionFeed init")
}
func (rp TimeSectionFeed) Run_interpose(request *data_type.Request) int {
	defer common.Timer("TimeSectionFeed", &(request.Timer_log), time.Now())
	day1List := make([]data_type.LiveInfo, 0)
	day2List := make([]data_type.LiveInfo, 0)
	day3List := make([]data_type.LiveInfo, 0)
	dayLess7List := make([]data_type.LiveInfo, 0)
	leftList := make([]data_type.LiveInfo, 0)
	for _, info := range request.Livelist {
		if time.Since(GetTime(info.LiveId)) < 24*time.Hour {
			day1List = append(day1List, info)
		} else if time.Since(GetTime(info.LiveId)) < 2*24*time.Hour {
			day2List = append(day2List, info)
		} else if time.Since(GetTime(info.LiveId)) < 3*24*time.Hour {
			day3List = append(day3List, info)
		} else if time.Since(GetTime(info.LiveId)) < 7*24*time.Hour {
			dayLess7List = append(dayLess7List, info)
		} else {
			leftList = append(leftList, info)
		}

	}
	day1List = append(day1List, day2List...)
	day1List = append(day1List, day3List...)
	day1List = append(day1List, dayLess7List...)
	day1List = append(day1List, leftList...)
	request.Livelist = day1List
	return 0
}
func GetTime(feedID string) time.Time {
	if len(feedID) < 10 {
		return time.Time{}
	}
	return common.TransTimestamp(feedID[:10])
}
