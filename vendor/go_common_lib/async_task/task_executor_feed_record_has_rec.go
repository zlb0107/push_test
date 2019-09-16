package async_task

import (
	"bytes"
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/prepare"
)

type ExecutorFeedRecordHasRec struct {
}

func init() {
	executor := NewExecutorFeedRecordHasRec()
	go executor.Run(ReadyQueue_feed_record_has_rec)
	logs.Warn("init ExecutorFeedRecordHasRec and run an ExecutorFeedRecordHasRec")
}

func NewExecutorFeedRecordHasRec() *ExecutorFeedRecordHasRec {
	return &ExecutorFeedRecordHasRec{}
}

func (this *ExecutorFeedRecordHasRec) Exec_task(task Task) bool {
	request := (*data_type.Request)(task)
	key := prepare.Gen_redis_key_has_rec(request)
	if key == "" {
		return false
	}
	return this.push_has_rec_list(key, request)
}

func (this *ExecutorFeedRecordHasRec) gen_has_rec_list_string(request *data_type.Request) string {
	//assert(request != nil)
	//we use bytes.Buffer
	count := 720 - len(request.Livelist)
	const split_char = ";"
	var buffer bytes.Buffer
	for k, v := range request.Has_rec_list {
		if k != "" && v {
			buffer.WriteString(k)
			buffer.WriteString(split_char)

			count--
			if count == 0 {
				break
			}
		}
	}
	for _, live_info := range request.Livelist {
		if live_info.LiveId != "" {
			buffer.WriteString(live_info.LiveId)
			buffer.WriteString(split_char)
		}
	}
	return buffer.String()
}
func (this *ExecutorFeedRecordHasRec) getExpire(request *data_type.Request) int {
	return 600
}
func (this *ExecutorFeedRecordHasRec) push_has_rec_list(key string, request *data_type.Request) bool {
	rc := prepare.Has_rec_redis.Get()
	defer rc.Close()
	new_str := this.gen_has_rec_list_string(request)
	expire := this.getExpire(request)
	//expire 10 minute
	_, err := rc.Do("setex", key, expire, new_str)
	if err != nil {
		logs.Error("setex error:", err, " key:", key)
		return false
	}
	return true
}

func (this *ExecutorFeedRecordHasRec) Run(listen_queue chan Task) {

	for {
		task := <-listen_queue
		if !this.Exec_task(task) {
			//logs.Error("exec failed")
		}
	}

}
