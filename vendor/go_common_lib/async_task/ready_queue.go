package async_task

import (
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"time"
)

type ReadyQueue chan Task

var (
	ReadyQueue_record_has_rec      ReadyQueue
	ReadyQueue_feed_record_has_rec ReadyQueue
)

const (
	default_queue_capacity = 1000
)

func init() {
	ReadyQueue_record_has_rec = make(chan Task, default_queue_capacity)
	ReadyQueue_feed_record_has_rec = make(chan Task, default_queue_capacity)
}

func AddAsyncTask(req *data_type.Request, queue ReadyQueue) bool {
	start := time.Now()
	if req == nil {
		return false
	}
	queue <- Task(req)
	common.Timer("AddAsyncTask", &(req.Timer_log), start)
	return true
}
