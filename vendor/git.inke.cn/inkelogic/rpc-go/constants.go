package rpc

import (
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys"
)

var (
	DISCOVERY_FILE = "./.discovery"
	APP_FILE       = "./.app"
)

var (
	SERVICE_REMOTE_INFO_PRE    = "/service_info/"
	SERVICE_CONFIG_CLIENT_TOML = "/config.toml"
)

var (
	ACCESS_LOG    string = "access"
	BUSSINESS_LOG string = "bussiness"
	SERVER_LOG    string = "server"
	STAT_LOG      string = "stat"
)

var (
	RPC_TYPE  = "rpc"
	HTTP_TYPE = "http"
)

var (
	LOG_ROTATE_HOUR  = "hour"
	LOG_ROTATE_DAY   = "day"
	LOG_ROTATE_MONTH = "month"
)

type ServerClient = daenerys.ServerClient

type JSONDataLogOption = daenerys.JSONDataLogOption

type RpcDefaultConfig struct {
	Server struct {
		ServiceName string   `toml:"service_name"`
		Port        int      `toml:"port"`
		Tags        []string `toml:"server_tags"`
		TCP         struct {
			IdleTimeout      int `toml:"idle_timeout"`
			KeepliveInterval int `toml:"keeplive_interval"`
		} `toml:"tcp"`
		HTTP struct {
			Location    string `toml:"location"`
			LogResponse string `toml:"logResponse"`
		} `toml:"http"`
	} `toml:"server"`
	Trace struct {
		Port int `toml:"port"`
	} `toml:"trace"`
	Monitor struct {
		AliveInterval int `toml:"alive_interval"`
	} `toml:"monitor"`
	Log struct {
		Level             string `toml:"level"`
		Rotate            string `toml:"rotate"`
		AccessRotate      string `toml:"access_rotate"`
		Accesslog         string `toml:"accesslog"`
		Businesslog       string `toml:"businesslog"`
		Serverlog         string `toml:"serverlog"`
		StatLog           string `toml:"statlog"`
		ErrorLog          string `toml:"errlog"`
		LogPath           string `toml:"logpath"`
		BalanceLogLevel   string `toml:"balance_log_level"`
		GenLogLevel       string `toml:"gen_log_level"`
		BusinessLogOff    bool   `toml:"business_log_off"`
		RequestBodyLogOff bool   `toml:"request_log_off"`
		SuccStatCode      []int  `toml:"succ_stat_code"`
	} `toml:"log"`
	ServerClient []ServerClient `toml:"server_client"`

	KafkaConsume []struct {
		ConsumeFrom    string `toml:"consume_from"`
		Zookeeperhost  string `toml:"zkpoints"`
		Topic          string `toml:"topic"`
		Group          string `toml:"group"`
		Initoffset     int    `toml:"initoffset"`
		ProcessTimeout int    `toml:"process_timeout"`
		CommitInterval int    `toml:"commit_interval"`
		GetError       bool   `toml:"get_error"`
		TraceEnable    bool   `toml:"trace_enable"`
	} `toml:"kafka_consume"`

	KafkaProducerClient []struct {
		ProducerTo   string `toml:"producer_to"`
		Broken       string `toml:"kafka_broken"`
		Retrymax     int    `toml:"retrymax"`
		RequiredAcks string `toml:"required_acks"`
		GetError     bool   `toml:"get_error"`
		GetSuccess   bool   `toml:"get_success"`
		UseSync      bool   `toml:"use_sync"`
	} `toml:"kafka_producer_client"`

	Redis []struct {
		ServerName     string `toml:"server_name"`
		Addr           string `toml:"addr"`
		Password       string `toml:"password"`
		MaxIdle        int    `toml:"max_idle"`
		MaxActive      int    `toml:"max_active"`
		IdleTimeout    int    `toml:"idle_timeout"`
		ConnectTimeout int    `toml:"connect_timeout"`
		ReadTimeout    int    `toml:"read_timeout"`
		WriteTimeout   int    `toml:"write_timeout"`
		Database       int    `toml:"database"`
		Retry          int    `toml:"retry"`
	} `toml:"redis"`

	Database []sql.SQLGroupConfig `toml:"database"`
	DataLog  JSONDataLogOption    `toml:"data_log"`
	Default  struct {
		HTTPMaxIdleConns        int `toml:"http_max_idle_conns"`
		HTTPMaxIdleConnsPerHost int `toml:"http_max_idle_conns_perhost"`
	} `toml:"default"`

	Circuit []daenerys.CircuitConfig `toml:"circuit"`
}
