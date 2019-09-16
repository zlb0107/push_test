package common

import (
	"strconv"
	"time"
)

func Timer(fun_name string, log **string, start time.Time) {
	dis := time.Since(start).Nanoseconds()
	msecond := dis / 1000000
	if msecond != 0 {
		old_log := ""
		if (*log) != nil {
			old_log = **log
		}
		new_log := old_log + " " + fun_name + ":" + strconv.FormatInt(msecond, 10)
		*log = &new_log
	}
}
func TimerV2(fun_name string, log **string, start time.Time, is_hit *bool) {
	dis := time.Since(start).Nanoseconds()
	msecond := dis / 1000000

	if msecond == 0 && *is_hit == false {
		return
	}
	old_log := ""
	if (*log) != nil {
		old_log = **log
	}
	hit_log := ""
	if *is_hit {
		hit_log = "_" + format_bool(*is_hit)
	}
	new_log := old_log + " " + fun_name + ":" + strconv.FormatInt(msecond, 10) + hit_log
	*log = &new_log
}
func format_bool(flag bool) string {
	if flag {
		return "1"
	} else {
		return "0"
	}
}
