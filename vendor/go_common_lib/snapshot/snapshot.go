package snapshot

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/kafka"
	"math/rand"
	"time"
)

var Channel chan string
var ChannelAfterVersion chan string
var ChannelPB chan string
var ChannelTrigger chan string

const Len_channel int = 10000
const KafkaAddress = "localkafka01.dsj.inke.srv:9092,localkafka02.dsj.inke.srv:9092,localkafka03.dsj.inke.srv:9092,localkafka04.dsj.inke.srv:9092,localkafka05.dsj.inke.srv:9092,localkafka06.dsj.inke.srv:9092,localkafka07.dsj.inke.srv:9092"

func init() {
	Channel = make(chan string, Len_channel)
	ChannelPB = make(chan string, Len_channel)
	ChannelAfterVersion = make(chan string, Len_channel)
	ChannelTrigger = make(chan string, Len_channel)
	go PutOnlinePbNormalkafka()
	go PutKafka(KafkaAddress, "rec_live_hall_snapshoot_content_pb", &(ChannelPB))
	go PutKafka(KafkaAddress, "newserverlog_rechall_trigger_details", &(ChannelTrigger))
}
func PutOnlinePbNormalkafka() {
	const kafka_address = "localkafka01.dsj.inke.srv:9092,localkafka02.dsj.inke.srv:9092,localkafka03.dsj.inke.srv:9092,localkafka04.dsj.inke.srv:9092,localkafka05.dsj.inke.srv:9092,localkafka06.dsj.inke.srv:9092,localkafka07.dsj.inke.srv:9092"
	//const kafka_address = "localhost:9092"
	p := kafka.Producer_init(kafka_address)

	i := 0
	for {
		i += 1
		if i >= 10000 {
			i %= 10000
			logs.Error("snapshot channel len:", len(ChannelAfterVersion))
		}
		if p == nil {
			logs.Error("p is nil")
			p = kafka.Producer_init(kafka_address)
			continue
		}

		snap := <-ChannelAfterVersion

		time_chan := make(chan int, 1)
		go func() {
			kafka.Async_producer(p, "rec_live_hall_snapshoot_content_online_pb", snap)
			time_chan <- 1
		}()
		ret := Func_timeout(time_chan, 1000)
		if ret == false {
			if p.Close() != nil {
				logs.Error("close p failed")
			}
			logs.Error("create a new p")
			p = kafka.Producer_init(kafka_address)
		}
	}
}
func PutNormalkafka() {
	const kafka_address = "ali-c-dsj-kafka01.bj:9092,ali-c-dsj-kafka02.bj:9092,ali-c-dsj-kafka03.bj:9092,ali-a-dsj-kafka04.bj:9092,ali-a-dsj-kafka05.bj:9092,ali-a-dsj-kafka06.bj:9092,ali-a-dsj-kafka07.bj:9092"
	p := kafka.Producer_init(kafka_address)
	for {
		if p == nil {
			logs.Error("p is nil")
			p = kafka.Producer_init(kafka_address)
			continue
		}
		rand.Seed(time.Now().UnixNano())
		randnum := rand.Float64()
		if randnum < 0.0001 {
			logs.Error("snapshot channel len:", len(Channel))
		}

		snap := <-Channel

		time_chan := make(chan int, 1)
		go func() {
			kafka.Async_producer(p, "rec_live_hall_snapshoot_content", snap)
			time_chan <- 1
		}()
		ret := Func_timeout(time_chan, 1000)
		if ret == false {
			if p.Close() != nil {
				logs.Error("close p failed")
			}
			logs.Error("create a new p")
			p = kafka.Producer_init(kafka_address)
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
				return false
			}
		}
	}
}
func PutKafka(kafka_address, topic string, channel *chan string) {
	p := kafka.Producer_init(kafka_address)
	(*channel) = make(chan string, Len_channel)
	i := 0
	for {
		i += 1
		if i >= 10000 {
			i %= 10000
			logs.Error("snapshot channel len:", len(*channel))
		}
		if p == nil {
			logs.Error("p is nil")
			p = kafka.Producer_init(kafka_address)
			continue
		}

		snap := <-*channel

		time_chan := make(chan int, 1)
		go func() {
			kafka.Async_producer(p, topic, snap)
			time_chan <- 1
		}()
		ret := Func_timeout(time_chan, 1000)
		if ret == false {
			if p.Close() != nil {
				logs.Error("close p failed")
			}
			logs.Error("create a new p")
			p = kafka.Producer_init(kafka_address)
		}
	}
}
