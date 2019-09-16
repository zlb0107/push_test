package daenerys

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/utils"
	"git.inke.cn/inkelogic/daenerys/log"
	dutils "git.inke.cn/inkelogic/daenerys/utils"
	clusterconfig "git.inke.cn/tpc/inf/go-upstream/config"
	upstreamConfig "git.inke.cn/tpc/inf/go-upstream/config"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	"git.inke.cn/tpc/inf/go-upstream/registry/consul"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
	"golang.org/x/net/trace"
)

type noopCloser struct{}

func (n noopCloser) Close() error {
	return nil
}

func readFile(name string) string {
	if b, err := ioutil.ReadFile(name); err != nil {
		return ""
	} else {
		return string(b)
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

func (d *Daenerys) verifyConfig() {
	if len(d.ConsulAddr) == 0 {
		d.ConsulAddr = consulAddr
	}
	if len(d.TraceReportAddr) == 0 {
		d.TraceReportAddr = traceReportAddr
	}
	// Init registry
	registry.Default, _ = consul.NewBackend(&upstreamConfig.Consul{Addr: d.ConsulAddr, Scheme: "http", Logger: logging.Log(logging.GenLoggerName)})
}

func (d *Daenerys) loggerInit() {
	if len(d.config.Log.LogPath) == 0 {
		d.config.Log.LogPath = "logs"
	}
	d.LogDir = d.config.Log.LogPath

	// Init common logger
	logging.InitCommonLog(logging.CommonLogConfig{
		Pathprefix:      d.config.Log.LogPath,
		Rotate:          d.config.Log.Rotate,
		GenLogLevel:     d.config.Log.GenLogLevel,
		BalanceLogLevel: d.config.Log.BalanceLogLevel,
	})

	// upstream logger
	upstream.SetLogger(logging.Log(logging.BalanceLoggerName))

	// internal logger
	rotateType := d.config.Log.Rotate
	var blog, glog, alog, slog log.Logger
	if d.config.Log.AccessLogOff {
		alog = log.NoopLogger()
	} else {
		alog = log.Default(filepath.Join(d.LogDir, "access.log")).Rotate(rotateType).Logger()
	}
	if d.config.Log.BusinessLogOff {
		blog = log.NoopLogger()
	} else {
		blog = log.Default(filepath.Join(d.LogDir, "bussiness.log")).Rotate(rotateType).Logger()
	}
	glog = log.Default(filepath.Join(d.LogDir, "gen.log")).Logger()
	slog = log.Default(filepath.Join(d.LogDir, "slow.log")).Rotate(rotateType).Logger()
	if d.Kit == nil {
		d.Kit = log.NewKit(blog, glog, alog, slog)
	}
}

// 如果设置了app_name则用app_name+service_name,如果没有则保持原有逻辑用service_name
// 此处逻辑保证注册与获取时service_name是一致的
func (d *Daenerys) injectServerClient(sc ServerClient) {
	sName := sc.ServiceName
	if sc.APPName != nil {
		if len(*sc.APPName) > 0 && *sc.APPName != INKE {
			sName = dutils.MakeAppServiceName(*sc.APPName, sc.ServiceName)
		}
	}
	if _, ok := d.serverClientMap.Load(sName); ok {
		return
	}
	cluster := d.makeCluster(sName, sc)
	if err := d.Clusters.InitService(cluster); err != nil {
		panic(err)
	}
	sc.Cluster = cluster
	d.serverClientMap.Store(sName, sc)
}

func (d *Daenerys) makeCluster(sName string, sc ServerClient) clusterconfig.Cluster {
	cluster := clusterconfig.NewCluster()
	if sc.ProtoType == "" || sc.ProtoType == "rpc" {
		sc.ProtoType = "http"
	}
	cluster.Name = fmt.Sprintf("%s-%s", sName, sc.ProtoType)
	if sc.APPName == nil { // 原有逻辑:使用本地环境的app_name
		cluster.Name = fmt.Sprintf("%s-%s", dutils.MakeAppServiceName(d.App, sc.ServiceName), sc.ProtoType)
	}
	cluster.StaticEndpoints = sc.Ipport
	if len(sc.Ipport) != 0 {
		// add fallback port
		var fallbackPort = ""
		if sc.ProtoType == "http" {
			fallbackPort = ":80"
		} else if sc.ProtoType == "https" {
			fallbackPort = ":443"
		}
		staticIPPorts := strings.Split(sc.Ipport, ",")
		for i := range staticIPPorts {
			_, _, err := net.SplitHostPort(staticIPPorts[i])
			if err != nil {
				if strings.Contains(err.Error(), "missing port") {
					staticIPPorts[i] = staticIPPorts[i] + fallbackPort
				}
			}
		}
		cluster.StaticEndpoints = strings.Join(staticIPPorts, ",")
	}
	cluster.Proto = sc.ProtoType
	cluster.LBType = sc.Balancetype
	cluster.EndpointsFrom = sc.EndpointsFrom
	cluster.CheckInterval = sc.CheckInterval
	cluster.UnHealthyThreshold = sc.UnHealthyThreshold
	cluster.HealthyThreshold = sc.HealthyThreshold
	cluster.LBPanicThreshold = sc.LBPanicThreshold
	cluster.LBSubsetKeys = sc.LBSubsetKeys
	cluster.LBDefaultKeys = sc.LBDefaultKeys
	cluster.Detector.DetectInterval = sc.DetectInterval
	cluster.Detector.ConsecutiveError = sc.ConsecutiveError
	cluster.Detector.ConsecutiveConnectionError = sc.ConsecutiveConnectionError
	cluster.Detector.MaxEjectionPercent = sc.MaxEjectionPercent
	cluster.Detector.SuccessRateMinHosts = sc.SuccessRateMinHosts
	cluster.Detector.SuccessRateRequestVolume = sc.SuccessRateRequestVolume
	cluster.Detector.SuccessRateStdevFactor = sc.SuccessRateStdevFactor
	cluster.Datacenter = sc.DC
	return cluster
}

func (d *Daenerys) initStat() {
	if len(d.LogDir) == 0 {
		utils.SetStat(filepath.Join(d.config.Log.LogPath, "stat"), d.localAppServiceName)
	} else {
		utils.SetStat(filepath.Join(d.LogDir, "stat"), d.localAppServiceName)
	}

	succCodeMap := make(map[int]int)
	for _, v := range d.config.Log.SuccStatCode {
		succCodeMap[v] = 1
	}
	utils.AddSBatchuccCode(succCodeMap)
}

func (d *Daenerys) initTrace() {
	tracePort := d.config.Trace.Port
	if tracePort != 0 {
		trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
			return true, true
		}
		http.HandleFunc("/_admin/gc", func(w http.ResponseWriter, r *http.Request) {
			runtime.GC()
			fmt.Fprint(w, "success\n")
		})
		http.HandleFunc("/_admin/free", func(w http.ResponseWriter, r *http.Request) {
			debug.FreeOSMemory()
			fmt.Fprint(w, "success\n")
		})
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", tracePort), nil)
	}
}

func (d *Daenerys) initCommLog() {
	if d.config.Log.Rotate == LOG_ROTATE_DAY {
		logging.SetRotateByDay()
	} else {
		logging.SetRotateByHour()
	}
	if len(d.config.Log.LogPath) > 0 {
		logging.SetOutputPath(d.config.Log.LogPath)
	} else {
		logging.SetOutputPath(d.LogDir)
	}
	if len(d.config.Log.Level) > 0 {
		logging.SetLevelByString(d.config.Log.Level)
	} else {
		logging.SetLevelByString(d.LogLevel)
	}

	if d.config.DataLog.Path != "" {
		var rotate logging.RollingFormat
		switch d.config.DataLog.Rotate {
		case LOG_ROTATE_HOUR:
			rotate = logging.HourlyRolling
		case LOG_ROTATE_DAY:
			rotate = logging.DailyRolling
		case LOG_ROTATE_MONTH:
			rotate = logging.MinutelyRolling
		default:
			rotate = logging.DailyRolling
		}
		name := d.Name
		if d.config.DataLog.TaskName != "" {
			name = d.config.DataLog.TaskName
		}
		err := logging.InitDataWithKey(d.config.DataLog.Path, rotate, name)
		if err != nil {
			panic(err)
		}
	}
}

func (d *Daenerys) initMiddleware() error {
	d.initCommLog()
	d.initStat()
	d.initTrace()

	if err := d.InitKafkaProducer(d.kafkaProducerConfig()); err != nil {
		return err
	}
	if err := d.InitKafkaConsume(d.config.KafkaConsume); err != nil {
		return err
	}
	if err := d.InitRedisClient(d.redisConfig()); err != nil {
		return err
	}
	if err := d.InitSqlClient(d.config.Database); err != nil {
		return err
	}
	/*if err := settingCircuit(d.config.circuitConfig()); err != nil {
		return err
	}*/
	return nil
}

func (d *Daenerys) kafkaProducerConfig() []kafka.KafkaProductConfig {
	var producerConfig []kafka.KafkaProductConfig
	if len(d.config.KafkaProducerClient) == 0 {
		return producerConfig
	}
	for _, defaultConfig := range d.config.KafkaProducerClient {
		var kpc kafka.KafkaProductConfig
		kpc.Broken = defaultConfig.Broken
		kpc.RetryMax = defaultConfig.RetryMax
		kpc.ProducerTo = defaultConfig.ProducerTo
		kpc.RequiredAcks = defaultConfig.Required_Acks
		kpc.GetError = defaultConfig.GetError
		kpc.GetSuccess = defaultConfig.GetSuccess
		kpc.UseSync = defaultConfig.Use_Sync
		producerConfig = append(producerConfig, kpc)
	}
	return producerConfig
}

func (d *Daenerys) redisConfig() []redis.RedisConfig {
	var redisConfig []redis.RedisConfig

	if len(d.config.Redis) == 0 {
		return redisConfig
	}

	for _, defaultConfig := range d.config.Redis {
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
