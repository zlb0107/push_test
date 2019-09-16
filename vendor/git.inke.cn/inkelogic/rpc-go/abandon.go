package rpc

import (
	"context"
	"fmt"
	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys"
	"git.inke.cn/inkelogic/daenerys/rpc/client"
	"strings"
)

var EnableTracing = true

var logdirGlobal = "logs"

// func StartMonitor(c Config) {}

func InitFrameworkUtils(c Config) {
	// nothing to do, daenerys init do it.
}

// SetBusinessLogFile  set client/server access log file path
func SetBusinessLogFile(path string) {
	// todo
}

func SetAccessRotateType(rotate string) {
	// todo
}

func SetAccessLogFile(path string) {
	// todo
}

func SyncLog() {
	// todo
	// accessLogger.Sync()
	// businessLogger.Sync()
	logging.Sync()
}

func NewRemoteTomlConfig() (*ConfigToml, error) {
	return GetTomlConfig(), nil
}

func Unregister() {
	daenerys.Default.Manager.Deregister()
}

func SetServiceName(name string) {
	daenerys.Default.Name = name
}

func InitRedisClient(redisConfigs []redis.RedisConfig) error {
	return daenerys.Default.InitRedisClient(redisConfigs)
}

func InitKafkaConsume(consumeConfigs []kafka.KafkaConsumeConfig) error {
	return daenerys.Default.InitKafkaConsume(consumeConfigs)
}

func InitKafkaProducer(producerConfigs []kafka.KafkaProductConfig) error {
	return daenerys.Default.InitKafkaProducer(producerConfigs)
}

func InitSQLClient(sqlConfig []sql.SQLGroupConfig) error {
	return daenerys.Default.InitSqlClient(sqlConfig)
}

func CloseAllClient() error {
	// should be delete
	return nil
}

func InitLoadBalance(serverClientList []ServerClient) error {
	for _, sc := range serverClientList {
		daenerys.Default.InjectServerClient(sc)
	}
	return nil
}

func GetRedis(service string) (r *redis.Redis, err error) {
	defer func() {
		if e := recover(); e != nil {
			r, err = nil, fmt.Errorf("%s", e)
		}
	}()
	return daenerys.RedisClient(context.TODO(), service), nil
}

func GetKafkaConsumeClient(consumeFrom string) (r *kafka.KafkaConsumeClient, err error) {
	defer func() {
		if e := recover(); e != nil {
			r, err = nil, fmt.Errorf("%s", e)
		}
	}()
	return daenerys.KafkaConsumeClient(context.TODO(), consumeFrom), nil
}

func GetSyncProducerClient(producerTo string) (r *kafka.KafkaSyncClient, err error) {
	defer func() {
		if e := recover(); e != nil {
			r, err = nil, fmt.Errorf("%s", e)
		}
	}()
	return daenerys.SyncProducerClient(context.TODO(), producerTo), nil
}

func GetKafkaProducerClient(producerTo string) (r *kafka.KafkaClient, err error) {
	defer func() {
		if e := recover(); e != nil {
			r, err = nil, fmt.Errorf("%s", e)
		}
	}()
	return daenerys.KafkaProducerClient(context.TODO(), producerTo), nil
}

func Invoke(ctx context.Context, method string, args, reply interface{}, cc *ClientConn) error {
	if value, ok := cc.endpoints.Load(method); !ok {
		cc.mu.Lock()

		if value, ok := cc.endpoints.Load(method); ok {
			cc.mu.Unlock()
			return value.(client.Client).Invoke(ctx, args, reply)
		}

		c := cc.factory.Client(method)
		cc.endpoints.Store(method, c)

		cc.mu.Unlock()
		return c.Invoke(ctx, args, reply)
	} else {
		return value.(client.Client).Invoke(ctx, args, reply)
	}
}

func getServiceTags(tags []string) map[string]string {
	serviceTags := make(map[string]string)
	for _, t := range tags {
		kvs := strings.SplitN(t, "=", 2)
		if len(kvs) > 1 {
			serviceTags[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
		}
	}
	return serviceTags
}
