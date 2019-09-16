package rpc

import (
	"bufio"
	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	serverLogPathName   = "log.server.path"
	businessLogPathName = "log.business.path"
	accessLogPathName   = "log.access.path"
	statLogPathName     = "stat.path"

	serverLogLevelName    = "log.server.level"
	defaultServerLogLevel = "debug"

	httpServeLocationName = "http.serve.location"
	httpServeLogBodyName  = "http.serve.log_response"

	portName   = "server.port"
	metricName = "stat.metric"

	idleTimeoutName       = "server.tcp.idle_timeout"
	keepAliveIntervalName = "server.tcp.keepalive_interval"
	NegotiateTimeoutName  = "server.tcp.negotiate_timeout"

	defaultServerLogPath   = "./server.log"
	defaultBusinessLogPath = "./business.log"
	defaultStatLogPath     = "./stat.log"
	defaultAccessLogPath   = "./access.log"

	defaultServerLocation = "/api/serve"
	defaultServerLogBody  = "false"

	tracePortName = "trace.server.port"

	defaultIdleTimeout       = 120
	defaultKeepaliveInterval = 10
	defaultNegotiateTimeout  = 5000
	alive_interval           = 10
)

type ConfigIni struct {
	prop map[string]string
}

func NewConfig(path string) (*ConfigIni, error) {
	daenerys.Init(daenerys.ConfigPath(path))
	fileObj, err := os.Open(path)
	if err != nil || fileObj == nil {
		return nil, err
	}
	br := bufio.NewReader(fileObj)
	prop := map[string]string{}
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		one := strings.TrimSpace(string(line))
		if strings.HasPrefix(one, "#") {
			continue
		}
		sect := strings.SplitN(one, "=", 2)
		if len(sect) != 2 {
			continue
		}
		prop[strings.TrimSpace(sect[0])] = strings.TrimSpace(sect[1])
	}
	toml = &ConfigIni{prop: prop}
	serviceName := toml.GetServiceName()
	if len(serviceName) == 0 {
		binaryName := strings.Split(os.Args[0], "/")
		serviceName = binaryName[len(binaryName)-1]
	}
	return toml.(*ConfigIni), nil
}

func (c *ConfigIni) Port() int {
	return c.IntWithDefault(portName, 4041)
}

func (c *ConfigIni) String(key string) (string, bool) {
	section, ok := c.prop[key]
	return section, ok
}

func (c *ConfigIni) StringWithDefault(key, defaultVal string) string {
	section, ok := c.prop[key]
	if !ok || section == "" {
		return defaultVal
	}
	return section
}

func (c *ConfigIni) Int(key string) (int, bool) {
	section, ok := c.prop[key]
	if !ok {
		return 0, false
	}
	data, err := strconv.Atoi(section)
	if err != nil {
		return 0, false
	}
	return data, true
}

func (c *ConfigIni) IntWithDefault(key string, value int) int {
	section, ok := c.prop[key]
	if !ok {
		return value
	}
	data, err := strconv.Atoi(section)
	if err != nil {
		return value
	}
	return data
}

func (c *ConfigIni) Set(key, value string) {
	c.prop[key] = value
}

func (c *ConfigIni) ServerLogPath() string {
	return c.StringWithDefault(serverLogPathName, defaultServerLogPath)
}

func (c *ConfigIni) ServerLogLevel() string {
	return c.StringWithDefault(serverLogLevelName, defaultServerLogLevel)
}
func (c *ConfigIni) BusinessLogPath() string {
	return c.StringWithDefault(businessLogPathName, defaultBusinessLogPath)
}

func (c *ConfigIni) StatLogPath() string {
	return c.StringWithDefault(statLogPathName, defaultStatLogPath)
}

func (c *ConfigIni) StatMetricName() string {
	return c.StringWithDefault(metricName, "")
}

func (c *ConfigIni) IdleTimeout() time.Duration {
	return time.Duration(c.IntWithDefault(idleTimeoutName, defaultIdleTimeout)) * time.Second
}

func (c *ConfigIni) KeepAliveInterval() time.Duration {
	return time.Duration(c.IntWithDefault(keepAliveIntervalName, defaultKeepaliveInterval)) * time.Second
}

func (c *ConfigIni) NegotiateTimeout() time.Duration {
	return time.Duration(c.IntWithDefault(NegotiateTimeoutName, defaultNegotiateTimeout)) * time.Millisecond
}

func (c *ConfigIni) TracePort() int {
	return c.IntWithDefault(tracePortName, 0)
}

func (c *ConfigIni) AccessLogPath() string {
	return c.StringWithDefault(accessLogPathName, defaultAccessLogPath)
}

func (c *ConfigIni) HTTPServeLocation() string {
	return c.StringWithDefault(httpServeLocationName, defaultServerLocation)
}

func (c *ConfigIni) HTTPServeLogBody() string {
	return c.StringWithDefault(httpServeLogBodyName, defaultServerLogBody)
}

func (tc *ConfigIni) GetServiceClients() []ServerClient {
	var clients []ServerClient
	return clients
}

func (tc *ConfigIni) GetKafkaProducerConfig() []kafka.KafkaProductConfig {
	var producerConfig []kafka.KafkaProductConfig
	return producerConfig
}

func (tc *ConfigIni) GetKafkaConsumeConfig() []kafka.KafkaConsumeConfig {
	var consumeConfig []kafka.KafkaConsumeConfig
	return consumeConfig
}

func (tc *ConfigIni) GetRedisConfig() []redis.RedisConfig {
	var redisConfig []redis.RedisConfig
	return redisConfig
}

func (tc *ConfigIni) GetSQLConfig() []sql.SQLGroupConfig {
	var c []sql.SQLGroupConfig
	return c
}

func (c *ConfigIni) GetServiceName() string {
	return c.StringWithDefault(metricName, "")
}

func (c *ConfigIni) GetMonitorInterval() int {

	return alive_interval
}

func (tc *ConfigIni) Rotate() string {
	return LOG_ROTATE_DAY
}

func (tc *ConfigIni) AccessRotate() string {
	return LOG_ROTATE_DAY
}

func (tc ConfigIni) LogPath() string {
	return ""
}

func (tc *ConfigIni) GetJSONDataLogOption() *JSONDataLogOption {
	return &JSONDataLogOption{Path: "./logs/trans.log", Rotate: "day"}
}

func (tc *ConfigIni) BusinessLogOff() bool {
	return false
}

func (tc *ConfigIni) RequestBodyLogOff() bool {
	return false
}

func (tc *ConfigIni) Tags() []string {
	return nil
}

func (tc *ConfigIni) SuccStatCode() []int {
	return nil
}

func (tc *ConfigIni) Circuit() []daenerys.CircuitConfig {
	return nil
}
