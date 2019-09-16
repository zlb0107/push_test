package daenerys

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/BackendPlatform/jaeger-client-go"

	jaegerconfig "git.inke.cn/BackendPlatform/jaeger-client-go/config"
	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/config"
	"git.inke.cn/inkelogic/daenerys/log"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	dutils "git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
	"github.com/opentracing/opentracing-go"
)

type Daenerys struct {
	// RunMode TODO
	RunMode Mode

	// Name is discovery name, it is from deploy platform by default.
	// Name will be used to register to discovery service
	Name string

	Namespace string

	JobName string

	App string

	Sync bool

	Version string

	Deps string

	LogDir string

	LogLevel string

	LogRotate string

	ConfigPath string

	Clusters *upstream.ClusterManager

	Manager *registry.ServiceManager

	// Tracer
	Tracer opentracing.Tracer

	Kit log.Kit

	// trace closer
	traceCloser io.Closer

	// config struct
	config daenerysConfig

	//config instance
	configInstance config.Config

	redisClientMap        map[string]*redis.Redis
	sqlClientMap          map[string]*sql.Group
	consumeClientMap      map[string]*kafka.KafkaConsumeClient
	producerClientMap     map[string]*kafka.KafkaClient
	syncProducerClientMap map[string]*kafka.KafkaSyncClient
	serverClientMap       *sync.Map
	httpClientMap         *sync.Map
	namespaceConfig       sync.Map
	rpcClientMap          sync.Map
	mu                    sync.Mutex
	localAppServiceName   string

	//Consul addr
	ConsulAddr string
	// Trace Addr
	TraceReportAddr string

	// pendingServiceInitTask
	pendingServerClientTask     []ServerClient
	pendingServerClientLock     sync.Mutex
	pendingServerClientTaskDone int32
	initOnce                    sync.Once
}

func New() *Daenerys {
	return &Daenerys{
		Sync:                        false,
		Name:                        "",
		Namespace:                   "",
		App:                         "",
		Version:                     "",
		RunMode:                     Production,
		LogDir:                      "logs",
		LogLevel:                    "debug",
		LogRotate:                   "hour",
		ConfigPath:                  "config.toml",
		traceCloser:                 noopCloser{},
		Clusters:                    upstream.NewClusterManager(),
		Manager:                     nil,
		redisClientMap:              make(map[string]*redis.Redis),
		sqlClientMap:                make(map[string]*sql.Group),
		consumeClientMap:            make(map[string]*kafka.KafkaConsumeClient),
		producerClientMap:           make(map[string]*kafka.KafkaClient),
		syncProducerClientMap:       make(map[string]*kafka.KafkaSyncClient),
		serverClientMap:             &sync.Map{},
		httpClientMap:               &sync.Map{},
		configInstance:              nil,
		pendingServerClientTask:     nil,
		pendingServerClientLock:     sync.Mutex{},
		pendingServerClientTaskDone: 0,
		initOnce:                    sync.Once{},
	}
}

func (d *Daenerys) Init(options ...Option) {
	d.initOnce.Do(func() {
		for _, opt := range options {
			opt(d)
		}

		if len(d.Name) == 0 {
			d.Name = strings.TrimSpace(readFile(".discovery"))
		}

		if len(d.JobName) == 0 {
			d.JobName = strings.TrimSpace(readFile(".jobname"))
		}

		if len(d.Namespace) == 0 && len(d.JobName) != 0 {
			// use cluster name as namespace name
			d.Namespace = d.JobName[strings.LastIndex(d.JobName, ".")+1:]
		}

		if len(d.Deps) == 0 {
			d.Deps = readFile(".deps")
		}

		if len(d.App) == 0 {
			d.App = strings.TrimSpace(readFile(".app"))
		}

		if len(d.Version) == 0 {
			d.Version = strings.TrimSpace(readFile(".version"))
		}

		c := d.Config().GetD("config.toml", d.ConfigPath)

		if err := c.Scan(&d.config); err != nil {
			panic(err)
		}

		// name must not empty.
		if len(d.Name) == 0 {
			//use binary name from command line
			if len(d.config.Server.ServiceName) != 0 {
				d.Name = d.config.Server.ServiceName
			} else if len(os.Args) > 0 {
				d.Name = filepath.Base(os.Args[0])
			}
		}

		d.configInstance = c
		if d.Sync {
			go func() {
				for {
					time.Sleep(time.Second * 5)
					if err := c.Sync(); err != nil {
						d.Kit.G().Log("message", "sync config from consul", "err", err)
					}
				}
			}()
		}

		d.loggerInit()

		d.verifyConfig()

		if d.Manager == nil {
			d.Manager = registry.NewServiceManager(logging.Log(logging.BalanceLoggerName))
		}

		// check start up config
		if d.RunMode == Development {
			if d.Kit == nil {
				d.Kit = log.NewKit(log.Stdout(), log.Stdout(), log.Stdout(), log.Stdout())
			}
			if d.Tracer == nil {
				d.Tracer = opentracing.NoopTracer{}
			}
		}

		d.localAppServiceName = dutils.MakeAppServiceName(d.App, d.Name)

		if d.Tracer == nil {
			// init tracer
			cfg := jaegerconfig.Configuration{
				// SamplingServerURL: "http://localhost:5778/sampling"
				Sampler: &jaegerconfig.SamplerConfig{Type: jaeger.SamplerTypeRemote},
				Reporter: &jaegerconfig.ReporterConfig{
					LogSpans:            false,
					BufferFlushInterval: 1 * time.Second,
					LocalAgentHostPort:  d.TraceReportAddr,
				},
			}
			tracer, closer, err := cfg.New(d.Name)
			if err != nil {
				panic(err)
			}
			d.traceCloser = closer
			d.Tracer = tracer
		}

		if d.Tracer != nil {
			opentracing.SetGlobalTracer(d.Tracer)
		}
		d.pendingServerClientLock.Lock()
		pending := d.pendingServerClientTask
		pending = append(pending, d.config.ServerClient...)
		d.pendingServerClientTask = nil
		atomic.StoreInt32(&d.pendingServerClientTaskDone, 1)
		d.pendingServerClientLock.Unlock()

		for _, sc := range pending {
			d.injectServerClient(sc)
		}
		if err := d.initMiddleware(); err != nil {
			panic(err)
		}
		fmt.Printf(
			"init daenerys success app:%s name:%s namespace:%s config:%s\n",
			d.App, d.Name, d.Namespace, d.ConfigPath,
		)

		breaker.InitWatcher(getRegistryKVPath(d.Name), d.Kit.G())
		ratelimit.InitWatcher(getRegistryKVPath(d.Name), d.Kit.G())
	})
}
