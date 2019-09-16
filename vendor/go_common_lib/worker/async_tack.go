package worker

import (
	"time"

	logs "github.com/cihub/seelog"
)

type Worker interface {
	Working() error
}

func RunWithRetry(w Worker, retry int) {
	for ; retry > 0; retry-- {
		if err := w.Working(); err == nil {
			break
		} else {
			logs.Error(err)
		}
	}
}

// RunTask 运行任务，w为任务，delay为第一次执行延时时间，interval为周期间隔时间，retry为任务返回err非空时重试次数
func RunTask(w Worker, delay, interval time.Duration, retry int) {
	// 先用timer定时到第一次延时
	timer := time.NewTimer(delay)
forLabel:
	for {
		select {
		case <-timer.C:
			go RunWithRetry(w, retry)
			break forLabel
		}
	}
	// 再使用ticker周期间隔执行
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go RunWithRetry(w, retry)
		}
	}
}
