package rpc

import (
	"time"

	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys"
)

type Config interface {
	Port() int
	String(key string) (string, bool)
	StringWithDefault(key, defaultVal string) string

	Int(key string) (int, bool)

	NegotiateTimeout() time.Duration

	IntWithDefault(key string, value int) int

	ServerLogPath() string
	ServerLogLevel() string
	BusinessLogPath() string
	AccessLogPath() string
	Rotate() string
	AccessRotate() string
	LogPath() string

	StatLogPath() string
	StatMetricName() string

	IdleTimeout() time.Duration
	KeepAliveInterval() time.Duration

	TracePort() int

	HTTPServeLocation() string
	HTTPServeLogBody() string

	GetServiceClients() []ServerClient

	GetKafkaConsumeConfig() []kafka.KafkaConsumeConfig
	GetKafkaProducerConfig() []kafka.KafkaProductConfig

	GetRedisConfig() []redis.RedisConfig

	GetSQLConfig() []sql.SQLGroupConfig

	GetServiceName() string

	GetMonitorInterval() int

	GetJSONDataLogOption() *JSONDataLogOption

	BusinessLogOff() bool
	RequestBodyLogOff() bool
	Circuit() []daenerys.CircuitConfig
	Tags() []string
	SuccStatCode() []int
}
