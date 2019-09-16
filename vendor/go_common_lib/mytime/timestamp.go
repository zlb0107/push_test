package common

import (
	logs "github.com/cihub/seelog"
	"strconv"
	"time"
)

func TransTimestamp(timestamp string) time.Time {
	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		logs.Error("err:", err)
		return time.Time{}
	}
	return time.Unix(t, 0)
}
