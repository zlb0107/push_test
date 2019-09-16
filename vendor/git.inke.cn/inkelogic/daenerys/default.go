package daenerys

import (
	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys/config"
	httpclient "git.inke.cn/inkelogic/daenerys/http/client"
	httpserver "git.inke.cn/inkelogic/daenerys/http/server"
	rpcclient "git.inke.cn/inkelogic/daenerys/rpc/client"
	rpcserver "git.inke.cn/inkelogic/daenerys/rpc/server"
	"golang.org/x/net/context"
)

var Default = New()

func Init(options ...Option) {
	Default.Init(options...)
}

func RPCServer() rpcserver.Server {
	return Default.RPCServer()
}

func HTTPServer() httpserver.Server {
	return Default.HTTPServer()
}

func Shutdown() error {
	return Default.Shutdown()
}

func Config() *config.Namespace {
	return Default.Config()
}

func ConfigInstance() config.Config {
	return Default.ConfigInstance()
}

func ConfigInstanceCtx(ctx context.Context) config.Config {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.ConfigInstance()
	}
	return Default.ConfigInstance()
}

func RPCFactory(ctx context.Context, name string) rpcclient.Factory {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.RPCFactory(name)
	}
	return Default.RPCFactory(name)
}

func HTTPClient(ctx context.Context, name string) httpclient.Client {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.HTTPClient(name)
	}
	return Default.HTTPClient(name)
}

func RedisClient(ctx context.Context, name string) *redis.Redis {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.RedisClient(name)
	}
	return Default.RedisClient(name)
}

func SQLClient(ctx context.Context, name string) *sql.Group {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.SQLClient(name)
	}
	return Default.SQLClient(name)
}

func KafkaConsumeClient(ctx context.Context, consumeFrom string) *kafka.KafkaConsumeClient {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.KafkaConsumeClient(consumeFrom)
	}
	return Default.KafkaConsumeClient(consumeFrom)
}

func KafkaProducerClient(ctx context.Context, consumeFrom string) *kafka.KafkaClient {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.KafkaProducerClient(consumeFrom)
	}
	return Default.KafkaProducerClient(consumeFrom)
}

func SyncProducerClient(ctx context.Context, producerTo string) *kafka.KafkaSyncClient {
	daenerys, ok := FromContext(ctx)
	if ok {
		return daenerys.SyncProducerClient(producerTo)
	}
	return Default.SyncProducerClient(producerTo)
}
