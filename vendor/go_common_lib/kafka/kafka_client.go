package kafka

import (
	logs "github.com/cihub/seelog"
)

func PutKafka(kafka_address, topic string, channel *chan string, channelLen int) {
	p := Producer_init(kafka_address)
	(*channel) = make(chan string, channelLen)
	i := 0
	for {
		i += 1
		if i >= 10000 {
			i %= 10000
			logs.Error("snapshot channel len:", len(*channel))
		}
		if p == nil {
			logs.Error("p is nil")
			p = Producer_init(kafka_address)
			continue
		}

		snap := <-*channel

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
