package daenerys

import (
	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/sql"
	upstreamconfig "git.inke.cn/tpc/inf/go-upstream/config"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

const INKE = "inke"

type ServerClient struct {
	APPName         *string `toml:"app_name"`
	ServiceName     string  `toml:"service_name"`
	Ipport          string  `toml:"endpoints"`
	Balancetype     string  `toml:"balancetype"`
	ProtoType       string  `toml:"proto"`
	ConnectTimeout  int     `toml:"connnect_timeout"`
	Namespace       string  `toml:"namespace"`
	ReadTimeout     int     `toml:"read_timeout"`
	WriteTimeout    int     `toml:"write_timeout"`
	MaxIdleConns    int     `toml:"max_idleconn"`
	RetryTimes      int     `toml:"retry_times"`
	SlowTime        int     `toml:"slow_time"`
	EndpointsFrom   string  `toml:"endpoints_from"`
	ConsulName      string  `toml:"consul_name"`
	LoadBalanceStat bool    `toml:"loadbalance_stat"`
	DC              string  `toml:"dc,omitempty"`

	// checker config
	CheckInterval      upstreamconfig.Duration `toml:"check_interval"`
	UnHealthyThreshold uint32                  `toml:"check_unhealth_threshold"`
	HealthyThreshold   uint32                  `toml:"check_healthy_threshold"`

	// lb advance config
	LBPanicThreshold int        `toml:"lb_panic_threshold"`
	LBSubsetKeys     [][]string `toml:"lb_subset_selectors"`
	LBDefaultKeys    []string   `toml:"lb_default_keys"`

	// detector config
	DetectInterval             upstreamconfig.Duration `toml:"detect_interval"`
	BaseEjectionDuration       upstreamconfig.Duration `toml:"base_ejection_duration"`
	ConsecutiveError           uint64                  `toml:"consecutive_error"`
	ConsecutiveConnectionError uint64                  `toml:"consecutive_connect_error"`
	MaxEjectionPercent         uint64                  `toml:"max_ejection_percent"`
	SuccessRateMinHosts        uint64                  `toml:"success_rate_min_hosts"`
	SuccessRateRequestVolume   uint64                  `toml:"success_rate_request_volume"`
	SuccessRateStdevFactor     float64                 `toml:"success_rate_stdev_factor"`
	Cluster                    upstreamconfig.Cluster

	Ratelimit []struct {
		Resource string `toml:"resource"`
		Limit    int    `toml:"limit"`
	} `toml:"ratelimit"`

	Breaker []struct {
		Resource   string `toml:"resource"`
		Open       bool   `toml:"open"`
		MinSamples int    `toml:"minsamples"`
		EThreshold int    `toml:"error_percent_threshold"`
		CThreshold int    `toml:"connsecutive_error_threshold"`
	} `toml:"breaker"`
}

type daenerysConfig struct {
	Server struct {
		ServiceName string   `toml:"service_name"`
		Port        int      `toml:"port"`
		Tags        []string `toml:"server_tags"`

		TCP struct {
			IdleTimeout      int `toml:"idle_timeout"`
			KeepliveInterval int `toml:"keeplive_interval"`
		} `toml:"tcp"`

		HTTP struct {
			Location    string `toml:"location"`
			LogResponse string `toml:"logResponse"`
		} `toml:"http"`

		Ratelimit []struct {
			Resource string `toml:"resource"`
			Peer     string `toml:"peer"`
			Limit    int    `toml:"limit"`
		} `toml:"ratelimit"`

		Breaker []struct {
			Resource   string `toml:"resource"`
			Open       bool   `toml:"open"`
			MinSamples int    `toml:"minsamples"`
			EThreshold int    `toml:"error_percent_threshold"`
			CThreshold int    `toml:"connsecutive_error_threshold"`
		} `toml:"breaker"`
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
		AccessLogOff      bool   `toml:"access_log_off"`
		BusinessLogOff    bool   `toml:"business_log_off"`
		RequestBodyLogOff bool   `toml:"request_log_off"`
		SuccStatCode      []int  `toml:"succ_stat_code"`
	} `toml:"log"`

	ServerClient        []ServerClient             `toml:"server_client"`
	KafkaConsume        []kafka.KafkaConsumeConfig `toml:"kafka_consume"`
	KafkaProducerClient []kafkaProducerItem        `toml:"kafka_producer_client"`
	Redis               []redisConfig              `toml:"redis"`
	Database            []sql.SQLGroupConfig       `toml:"database"`
	Circuit             []CircuitConfig            `toml:"circuit"`
	DataLog             JSONDataLogOption          `toml:"data_log"`
}

type JSONDataLogOption struct {
	Path     string `toml:"path"`
	Rotate   string `toml:"rotate"`
	TaskName string `toml:"task_name"`
}

type CircuitConfig struct {
	Type       string   `toml:"type"`
	Service    string   `toml:"service"`
	Resource   string   `toml:"resource"`
	End        string   `toml:"end"`
	Open       bool     `toml:"open"`
	Threshold  float64  `toml:"threshold"`
	Strategy   string   `toml:"strategy"`
	MinSamples int64    `toml:"minsamples"`
	RT         duration `toml:"rt"`
}

type kafkaProducerItem struct {
	kafka.KafkaProductConfig
	Required_Acks string `toml:"required_acks"` //old rpc-go
	Use_Sync      bool   `toml:"use_sync"`      //old rpc-go
}

//golang包中的redis是json格式,此处转为toml格式
type redisConfig struct {
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
	SlowTime       int    `toml:"slow_time"`
	Retry          int    `toml:"retry"`
}

type ClientOption func(*ClientOptions)

type ClientOptions struct{}
