package kafka

import (
	"time"

	logs "github.com/cihub/seelog"
)

func Put_kafka(kafka_address, topic string, Channel *chan string, timeout int) {
	p := Producer_init(kafka_address)
	for {
		if p == nil {
			logs.Error("p is nil")
			p = Producer_init(kafka_address)
			continue
		}
		if time.Now().UnixNano()%10000 == 0 {
			logs.Error("snapshot channel len:", len(*Channel))
		}

		snap := <-*Channel
		//kafka.Async_producer(p, "rec_live_hall_snapshoot_content", snap)

		time_chan := make(chan int, 1)
		go func() {
			Async_producer(p, topic, snap)
			time_chan <- 1
		}()
		ret := Func_timeout(time_chan, 1000)
		if ret == false {
			if p.Close() != nil {
				logs.Error("close p failed")
			}
			logs.Error("create a new p")
			p = Producer_init(kafka_address)
		}
	}
}
func Func_timeout(time_chan chan int, timeout int64) bool {
	for {
		select {
		case <-time_chan:
			{
				return true
			}
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			{
				//close(time_chan)
				return false
			}
		}
	}
}
