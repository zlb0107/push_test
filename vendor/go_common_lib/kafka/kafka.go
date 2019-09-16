package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster" //support automatic consumer-group rebalancing and offset tracking
	logs "github.com/cihub/seelog"
	"github.com/sdbaiguanghe/glog"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	topics = "aliyun_kfk_applog_common"
)

// consumer 消费者
func Consumer(mytopics string, hosts string, groupID string, channel chan string) {
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始

	c, err := cluster.NewConsumer(strings.Split(hosts, ","), groupID, strings.Split(mytopics, ","), config)
	if err != nil {
		glog.Errorf("Failed open consumer: %v", err)
		return
	}
	defer c.Close()
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				glog.Errorln(err)
			case <-noti:
				{
					//					glog.Errorln("get some")
				}
				//			case <-time.After(time.Millisecond * 1000):
				//				{
				//					glog.Errorln("timeout")
				//				}
			}
		}
	}(c)

	for msg := range c.Messages() {
		channel <- string(msg.Value)
		//		fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
	}
}

// syncProducer 同步生产者
// 并发量小时，可以用这种方式
func syncProducer() {
	config := sarama.NewConfig()
	//  config.Producer.RequiredAcks = sarama.WaitForAll
	//  config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(strings.Split("localhost:9092", ","), config)
	defer p.Close()
	if err != nil {
		glog.Errorln(err)
		return
	}

	v := "sync: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	fmt.Fprintln(os.Stdout, v)
	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.ByteEncoder(v),
	}
	if _, _, err := p.SendMessage(msg); err != nil {
		glog.Errorln(err)
		return
	}
}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func Producer_init(brokers string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewAsyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		logs.Error("init producer failed:", err)
		return nil
	}
	return p
}
func Stack() []byte {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return buf[:n]
}
func Async_producer(p sarama.AsyncProducer, mytopic, v string) {
	// 在抛出异常时，打印堆栈信息
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Controller: Panic. panic message: %#v. stack info: \n%s", r, Stack())
			logs.Error(msg)
			logs.Flush()
		}
	}()

	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				{
					if err != nil {
						logs.Error("kafka:", err)
					}
					break
				}
			case <-success:
				{
					break
				}
			case <-time.After(time.Millisecond * 10):
				{
					//					logs.Error("timout:", v)
					break
				}
			}
			break
		}
	}(p)

	//	v := "async: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	//	fmt.Fprintln(os.Stdout, v)
	msg := &sarama.ProducerMessage{
		Topic: mytopic,
		Value: sarama.ByteEncoder(v),
	}
	p.Input() <- msg
}
