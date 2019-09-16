package rpc

import (
	_ "net/http/pprof"
	"strings"
	"time"

	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/BackendPlatform/golang/tomlconfig"
	"git.inke.cn/inkelogic/daenerys"
	ikconfig "git.inke.cn/inkelogic/daenerys/config"
)

type ConfigToml struct {
	config        *tomlconfig.Config
	defaultConfig *RpcDefaultConfig
}

var (
	logRotateType string
	toml          Config
	tomlDefault   *RpcDefaultConfig
)

func init() {
	logRotateType = LOG_ROTATE_HOUR
}

func GetServiceName() string {
	return getAppServiceName(daenerys.Default.App, daenerys.Default.Name)
}

func getAppServiceName(app, name string) string {
	if len(app) == 0 {
		return name
	}
	return app + "." + name
}

func RotateDay() bool {
	if logRotateType == LOG_ROTATE_DAY {
		return true
	}
	return false
}

func GetLogRotateType() string {
	return logRotateType
}

func makeLogPath(path string, logname string) string {
	logpath := path + "/" + logname
	return logpath
}

func GetFileConfig() *RpcDefaultConfig {
	return tomlDefault
}

func NewRemoteConfigToml(v interface{}) (*ConfigToml, error) {
	return NewConfigToml("", v)
}

func NewConfigToml(path string, v interface{}) (*ConfigToml, error) {
	daenerys.Init(daenerys.ConfigPath(path))
	if err := ikconfig.Files(path); err != nil {
		return nil, err
	}
	if err := ikconfig.Consul(ikconfig.DefaultRemotePath(GetServiceName(), SERVICE_CONFIG_CLIENT_TOML)); err != nil {
		logging.Infof("read from consul failed: %v\n", err)
	}
	if err := ikconfig.Scan(v); err != nil {
		return nil, err
	}
	tomlDefault = &RpcDefaultConfig{}
	if err := ikconfig.Scan(tomlDefault); err != nil {
		return nil, err
	}
	toml = &ConfigToml{
		defaultConfig: tomlDefault,
	}
	return toml.(*ConfigToml), nil
}

func GetTomlConfig() *ConfigToml {
	return toml.(*ConfigToml)
}

func (tc *ConfigToml) Port() int {
	return int(tc.defaultConfig.Server.Port)
}

func (tc *ConfigToml) SuccStatCode() []int {
	return tc.defaultConfig.Log.SuccStatCode
}

func (tc *ConfigToml) String(key string) (string, bool) {
	if len(key) == 0 {
		return "", false
	}
	subKeys := strings.SplitN(strings.Trim(strings.TrimSpace(key), "."), ".", -1)
	result := ikconfig.Get(subKeys...)
	if result == nil {
		return "", false
	}
	return result.String(""), true
}

func (tc *ConfigToml) StringWithDefault(key, defaultVal string) string {
	if len(key) == 0 {
		return defaultVal
	}
	subKeys := strings.SplitN(strings.Trim(strings.TrimSpace(key), "."), ".", -1)
	result := ikconfig.Get(subKeys...)
	if result != nil {
		return result.String(defaultVal)
	}
	return defaultVal
}

func (tc *ConfigToml) Int(key string) (int, bool) {
	if len(key) == 0 {
		return 0, false
	}
	subKeys := strings.SplitN(strings.Trim(strings.TrimSpace(key), "."), ".", -1)
	result := ikconfig.Get(subKeys...)
	if result == nil {
		return 0, false
	} else {
		return result.Int(0), true
	}
}

func (tc *ConfigToml) IntWithDefault(key string, defaultVal int) int {
	if len(key) == 0 {
		return 0
	}
	subKeys := strings.SplitN(strings.Trim(strings.TrimSpace(key), "."), ".", -1)
	result := ikconfig.Get(subKeys...)
	if result != nil {
		return result.Int(defaultVal)
	}
	return defaultVal
}

func (tc *ConfigToml) ServerLogPath() string {
	return tc.defaultConfig.Log.Serverlog
}

func (tc *ConfigToml) ServerLogLevel() string {
	return tc.defaultConfig.Log.Level
}

func (tc *ConfigToml) BusinessLogPath() string {
	return tc.defaultConfig.Log.Businesslog
}

func (tc *ConfigToml) StatLogPath() string {
	return tc.defaultConfig.Log.StatLog
}

func (tc *ConfigToml) Rotate() string {
	return tc.defaultConfig.Log.Rotate
}

func (tc *ConfigToml) AccessRotate() string {
	return tc.defaultConfig.Log.AccessRotate
}

func (tc ConfigToml) LogPath() string {
	return tc.defaultConfig.Log.LogPath
}

func (tc *ConfigToml) StatMetricName() string {
	return tc.defaultConfig.Server.ServiceName
}

func (tc *ConfigToml) IdleTimeout() time.Duration {
	return time.Duration(tc.defaultConfig.Server.TCP.IdleTimeout) * time.Second
}

func (tc *ConfigToml) KeepAliveInterval() time.Duration {
	return time.Duration(tc.defaultConfig.Server.TCP.KeepliveInterval) * time.Second
}

func (tc *ConfigToml) TracePort() int {
	return tc.defaultConfig.Trace.Port
}

func (tc *ConfigToml) AccessLogPath() string {
	return tc.defaultConfig.Log.Accesslog
}

func (tc *ConfigToml) HTTPServeLocation() string {
	return tc.defaultConfig.Server.HTTP.Location
}

func (tc *ConfigToml) HTTPServeLogBody() string {
	return tc.defaultConfig.Server.HTTP.LogResponse
}

func (tc *ConfigToml) NegotiateTimeout() time.Duration {
	return time.Duration(defaultNegotiateTimeout) * time.Second
}

func (tc *ConfigToml) GetServiceClients() []ServerClient {
	var clients []ServerClient
	for _, client := range tc.defaultConfig.ServerClient {
		var sc ServerClient
		sc = client
		// old logic: add appname prefix
		sc.ServiceName = getAppServiceName(daenerys.Default.App, client.ServiceName)
		clients = append(clients, sc)
	}
	return clients
}

func (tc *ConfigToml) GetKafkaConsumeConfig() []kafka.KafkaConsumeConfig {
	var consumeConfig []kafka.KafkaConsumeConfig
	if len(tc.defaultConfig.KafkaConsume) == 0 {
		return consumeConfig
	}
	for _, defaultConfig := range tc.defaultConfig.KafkaConsume {
		var kcc kafka.KafkaConsumeConfig
		kcc.CommitInterval = defaultConfig.CommitInterval
		kcc.ConsumeFrom = defaultConfig.ConsumeFrom
		kcc.Group = defaultConfig.Group
		kcc.Initoffset = defaultConfig.Initoffset
		kcc.Zookeeperhost = defaultConfig.Zookeeperhost
		kcc.ProcessTimeout = defaultConfig.ProcessTimeout
		kcc.Topic = defaultConfig.Topic
		kcc.GetError = defaultConfig.GetError
		kcc.TraceEnable = defaultConfig.TraceEnable
		consumeConfig = append(consumeConfig, kcc)
	}
	return consumeConfig
}

func (tc *ConfigToml) GetKafkaProducerConfig() []kafka.KafkaProductConfig {
	var producerConfig []kafka.KafkaProductConfig
	if len(tc.defaultConfig.KafkaProducerClient) == 0 {
		return producerConfig
	}
	for _, defaultConfig := range tc.defaultConfig.KafkaProducerClient {
		var kpc kafka.KafkaProductConfig
		kpc.Broken = defaultConfig.Broken
		kpc.RetryMax = defaultConfig.Retrymax
		kpc.ProducerTo = defaultConfig.ProducerTo
		kpc.RequiredAcks = defaultConfig.RequiredAcks
		kpc.GetError = defaultConfig.GetError
		kpc.GetSuccess = defaultConfig.GetSuccess
		kpc.UseSync = defaultConfig.UseSync
		producerConfig = append(producerConfig, kpc)
	}
	return producerConfig
}

func (tc *ConfigToml) GetRedisConfig() []redis.RedisConfig {
	var redisConfig []redis.RedisConfig
	if len(tc.defaultConfig.Redis) == 0 {
		return redisConfig
	}
	for _, defaultConfig := range tc.defaultConfig.Redis {
		var rc redis.RedisConfig
		rc.ServerName = defaultConfig.ServerName
		rc.Addr = defaultConfig.Addr
		rc.Password = defaultConfig.Password
		rc.MaxIdle = defaultConfig.MaxIdle
		rc.MaxActive = defaultConfig.MaxActive
		rc.IdleTimeout = defaultConfig.IdleTimeout
		rc.ConnectTimeout = defaultConfig.ConnectTimeout
		rc.ReadTimeout = defaultConfig.ReadTimeout
		rc.WriteTimeout = defaultConfig.WriteTimeout
		rc.Database = defaultConfig.Database
		rc.Retry = defaultConfig.Retry
		redisConfig = append(redisConfig, rc)
	}
	return redisConfig
}

func (tc *ConfigToml) GetSQLConfig() []sql.SQLGroupConfig {
	return tc.defaultConfig.Database
}

func (tc *ConfigToml) GetServiceName() string {
	return tc.defaultConfig.Server.ServiceName
}

func (tc *ConfigToml) GetMonitorInterval() int {
	timeInterval := tc.defaultConfig.Monitor.AliveInterval
	if timeInterval == 0 {
		timeInterval = 10
	}
	return timeInterval
}

func (tc *ConfigToml) GetJSONDataLogOption() *JSONDataLogOption {
	return &tc.defaultConfig.DataLog
}

func (tc *ConfigToml) BusinessLogOff() bool {
	return tc.defaultConfig.Log.BusinessLogOff
}

func (tc *ConfigToml) RequestBodyLogOff() bool {
	return tc.defaultConfig.Log.RequestBodyLogOff
}

func (tc *ConfigToml) Circuit() []daenerys.CircuitConfig {
	return tc.defaultConfig.Circuit
}

func (tc *ConfigToml) Tags() []string {
	return tc.defaultConfig.Server.Tags
}
