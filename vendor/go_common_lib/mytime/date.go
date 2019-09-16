package common

import (
	"strconv"
	"time"
)

func Get_date_timestamp() string {
	year, month, day := time.Now().Date()
	return strconv.Itoa(year) + "_" + strconv.Itoa(int(month)) + "_" + strconv.Itoa(day)
}
