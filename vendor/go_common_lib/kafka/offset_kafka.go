package kafka

import (
	log "github.com/cihub/seelog"
	"strings"
	//	"encoding/json"
	"gopkg.in/Shopify/sarama.v1"
	//	"go_common_lib/data_type"
)

//const TOPIC string = "aliyun_kfk_applog_action"
const TOPIC string = "user_feed_logic_feed"

//const ADDRS string = "hadoop102:9092,hadoop103:9092,hadoop104:9092,hadoop105:9092,hadoop106:9092,hadoop107:9092,hadoop108:9092"
const ADDRS string = "ali-c-dsj-kafka01.bj:9092,ali-c-dsj-kafka02.bj:9092,ali-c-dsj-kafka03.bj:9092,ali-a-dsj-kafka04.bj:9092,ali-a-dsj-kafka05.bj:9092,ali-a-dsj-kafka06.bj:9092,ali-a-dsj-kafka07.bj:9092"

func InitKafka(addrs, group, topic string, useNew bool) (sarama.Consumer, sarama.OffsetManager, []int32) {
	config := sarama.NewConfig()
	if useNew {
		config.Version = sarama.V0_11_0_0
	}
	temp := strings.Split(addrs, ",")
	client, err := sarama.NewClient(temp, config)
	if err != nil {
		log.Error(err)
		panic("kafka")
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		log.Error(err)
		panic("kafka")
	}

	pids, err := client.Partitions(topic)
	if err != nil {
		log.Error(err)
		panic("kafka")
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Error(err)
		panic("kafka")
	}

	return consumer, offsetManager, pids
}

func Consume(c sarama.Consumer, om sarama.OffsetManager, p int32, offset int64, contentChan *chan *sarama.ConsumerMessage, topic string) {
	pom, err := om.ManagePartition(topic, p)
	if err != nil {
		log.Error(err, " p:", p)
		return
	}
	defer pom.Close()

	//offset, _ := pom.NextOffset()
	if offset == -1 {
		offset = sarama.OffsetOldest
	}

	pc, err := c.ConsumePartition(topic, p, offset)
	if err != nil {
		log.Error("~~~~~~~~~~~~~", err, " topic:", topic, " p:", p, " offset:", offset)
		log.Flush()
		pc, _ = c.ConsumePartition(topic, p, sarama.OffsetOldest)
	} else {
		defer pc.Close()
	}

	for msg := range pc.Messages() {
		(*contentChan) <- msg
		//		log.Warn("[%v] Consumed message offset %v\n", p, msg.Offset)
		//		log.Warn("msg_k:", string(msg.Key))
		//	log.Warn("msg_v:", string(msg.Value), " len:", len(*contentChan))
		pom.MarkOffset(msg.Offset+1, "")
	}
}

func ConsumeExtend(c sarama.Consumer, om sarama.OffsetManager, p int32, offset int64, contentChan *chan *sarama.ConsumerMessage, topic string) {
	pom, err := om.ManagePartition(topic, p)
	if err != nil {
		log.Error(err, " p:", p)
		return
	}
	defer pom.Close()

	//offset, _ := pom.NextOffset()
	if offset == -1 {
		offset = sarama.OffsetOldest
	} else if offset == -2 {
		offset = sarama.OffsetNewest //实际值是-1
	} else if offset == -3 { //从之前的位置开始读，如果没有这个group的偏移量，则从最老的位置开始读
		offset, _ = pom.NextOffset()
		if offset == sarama.OffsetNewest {
			offset = sarama.OffsetOldest
		}
	}

	pc, err := c.ConsumePartition(topic, p, offset)
	if err != nil {
		log.Error("~~~~~~~~~~~~~", err, " topic:", topic, " p:", p, " offset:", offset)
		log.Flush()
		pc, _ = c.ConsumePartition(topic, p, sarama.OffsetOldest)
	} else {
		defer pc.Close()
	}

	for msg := range pc.Messages() {
		(*contentChan) <- msg
		//		log.Warn("[%v] Consumed message offset %v\n", p, msg.Offset)
		//		log.Warn("msg_k:", string(msg.Key))
		//	log.Warn("msg_v:", string(msg.Value), " len:", len(*contentChan))
		pom.MarkOffset(msg.Offset+1, "")
	}
}
