package daenerys

import (
	"fmt"
	"sync/atomic"
	"time"

	"git.inke.cn/BackendPlatform/golang/kafka"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/BackendPlatform/golang/redis"
	"git.inke.cn/BackendPlatform/golang/sql"
	"git.inke.cn/inkelogic/daenerys/breaker"
	"git.inke.cn/inkelogic/daenerys/config"
	httpclient "git.inke.cn/inkelogic/daenerys/http/client"
	httpserver "git.inke.cn/inkelogic/daenerys/http/server"
	"git.inke.cn/inkelogic/daenerys/ratelimit"
	rpcclient "git.inke.cn/inkelogic/daenerys/rpc/client"
	"git.inke.cn/inkelogic/daenerys/rpc/codec"
	rpcserver "git.inke.cn/inkelogic/daenerys/rpc/server"
	dutils "git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	"github.com/hashicorp/go-multierror"
)

func (d *Daenerys) Shutdown() error {
	var result error

	// will be blocked
	//for _, client := range d.consumeClientMap {
	//	result = multierror.Append(result, client.Close())
	//}
	for _, client := range d.producerClientMap {
		result = multierror.Append(result, client.Close())
	}
	for _, client := range d.syncProducerClientMap {
		result = multierror.Append(result, client.Close())
	}

	d.traceCloser.Close()
	logging.Sync()

	return result
}

func (d *Daenerys) Config() *config.Namespace {
	namespace := d.Namespace
	if n, ok := d.namespaceConfig.Load(namespace); ok {
		return n.(*config.Namespace)
	}

	n, _ := d.namespaceConfig.LoadOrStore(
		namespace,
		config.NewNamespace(getRegistryKVPath(d.Name)).With(namespace),
	)
	return n.(*config.Namespace)
}

func (d *Daenerys) ConfigInstance() config.Config {
	return d.configInstance
}

func (d *Daenerys) InjectServerClient(sc ServerClient) {
	if atomic.LoadInt32(&d.pendingServerClientTaskDone) == 0 {
		d.pendingServerClientLock.Lock()
		defer d.pendingServerClientLock.Unlock()
		if atomic.LoadInt32(&d.pendingServerClientTaskDone) == 0 {
			d.pendingServerClientTask = append(d.pendingServerClientTask, sc)
			return
		}
	}
	d.injectServerClient(sc)
}

// 对于以下方法中的service参数说明:
// 如果对应的server_client配置了app_name选项,则需要调用方保证service参数带上app_name前缀
// 如果没有配置,则保持原有逻辑,	service参数不用改动
func (d *Daenerys) FindServerClient(service string) (ServerClient, error) {
	if value, ok := d.serverClientMap.Load(service); ok {
		sc := value.(ServerClient)
		if sc.ReadTimeout == 0 {
			sc.ReadTimeout = int(200 * time.Millisecond)
		}
		if sc.ConnectTimeout == 0 {
			sc.ConnectTimeout = int(200 * time.Millisecond)
		}
		return sc, nil
	}
	return ServerClient{}, fmt.Errorf("client config for %s not exist", service)
}

func (d *Daenerys) ServiceClientWithApp(appName, serviceName string) (ServerClient, error) {
	appServiceName := dutils.MakeAppServiceName(appName, serviceName)
	return d.FindServerClient(appServiceName)
}

// RPC create a new rpc client instance, default use http protocol.
func (d *Daenerys) RPCFactory(name string) rpcclient.Factory {
	if c, ok := d.rpcClientMap.Load(name); ok {
		return c.(rpcclient.Factory)
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	if c, ok := d.rpcClientMap.Load(name); ok {
		return c.(rpcclient.Factory)
	}

	sc, err := d.FindServerClient(name)
	if err != nil {
		panic(fmt.Errorf("namespace %s rpcclient %s not exist: %s", d.Namespace, name, err))
	}
	sName := sc.ServiceName

	limitMaps := map[string]int{}
	for _, c := range sc.Ratelimit {
		limitMaps[c.Resource] = c.Limit
	}

	var clusterName string
	if sc.APPName != nil {
		clusterName = fmt.Sprintf("%s-http", sName)
	} else {
		clusterName = fmt.Sprintf("%s-http", dutils.MakeAppServiceName(d.App, sc.ServiceName))
	}

	client := rpcclient.HFactory(
		rpcclient.Cluster(d.Clusters.Cluster(clusterName)),
		rpcclient.Kit(d.Kit),
		rpcclient.Codec(codec.NewJSONCodec()),
		rpcclient.Tracer(d.Tracer),
		rpcclient.Retries(sc.RetryTimes),
		rpcclient.RequestTimeout(time.Duration(sc.ReadTimeout)*time.Millisecond),
		rpcclient.Name(sName),
		rpcclient.SDName(d.localAppServiceName),
		rpcclient.Namespace(sc.Namespace),
		rpcclient.Ratelimit(
			ratelimit.NewConfig(
				d.Namespace,
				name,
				limitMaps,
			),
		),
		rpcclient.Breaker(
			breaker.NewConfig(
				d.Namespace,
				name,
				getBreakerConfig(sc),
			),
		),
	)
	d.rpcClientMap.Store(name, client)
	return client
}

func (d *Daenerys) RPCServer() rpcserver.Server {
	port := d.config.Server.Port
	if port == 0 {
		panic("server port is 0")
	}

	limitMaps := map[string]int{}
	for _, c := range d.config.Server.Ratelimit {
		limitMaps[c.Resource] = c.Limit
	}

	server := rpcserver.BothServer(
		port,
		rpcserver.Tags(getServiceTags(d.config.Server.Tags)),
		rpcserver.Manager(d.Manager),
		rpcserver.LoggerKit(d.Kit),
		rpcserver.Tracer(d.Tracer),
		rpcserver.Name(d.localAppServiceName),
		rpcserver.Registry(registry.Default),
		rpcserver.Ratelimit(
			ratelimit.NewConfig(
				d.Namespace,
				"",
				limitMaps,
			),
		),
		rpcserver.Breaker(
			breaker.NewConfig(
				//d.Config().With("breaker").With("server"),
				d.Namespace,
				"",
				//d.Kit.G(),
				getBreakerConfigServer(d.config),
			),
		),
	)
	return server
}

func (d *Daenerys) HTTPClient(name string) httpclient.Client {

	sc, err := d.FindServerClient(name)
	if err != nil {
		panic(fmt.Errorf("namespace %s httpclient %s not exist: %s", d.Namespace, name, err))
	}
	if sc.ProtoType == "" || sc.ProtoType == "rpc" {
		sc.ProtoType = "http"
	}
	sName := sc.ServiceName
	if sc.APPName != nil {
		if len(*sc.APPName) > 0 && *sc.APPName != INKE {
			sName = fmt.Sprintf("%s.%s", *sc.APPName, sc.ServiceName)
		}
	}
	if v, ok := d.httpClientMap.Load(sName); ok {
		return v.(httpclient.Client)
	}
	var clusterName string
	if sc.APPName != nil {
		clusterName = fmt.Sprintf("%s-%s", sName, sc.ProtoType)
	} else {
		clusterName = fmt.Sprintf("%s-%s", dutils.MakeAppServiceName(d.App, sc.ServiceName), sc.ProtoType)
	}

	limitMaps := map[string]int{}
	for _, c := range sc.Ratelimit {
		limitMaps[c.Resource] = c.Limit
	}

	client := httpclient.NewClient(
		httpclient.Cluster(d.Clusters.Cluster(clusterName)),
		httpclient.Tracer(d.Tracer),
		httpclient.Logger(d.Kit),
		httpclient.MaxIdleConns(sc.MaxIdleConns),
		httpclient.DialTimeout(time.Duration(sc.ConnectTimeout)*time.Millisecond),
		httpclient.KeepAliveTimeout(30*time.Second),
		httpclient.RetryTimes(sc.RetryTimes),
		httpclient.Namespace(sc.Namespace),
		httpclient.Ratelimit(
			ratelimit.NewConfig(
				d.Namespace,
				name,
				limitMaps,
			),
		),
		httpclient.Breaker(
			breaker.NewConfig(
				d.Namespace,
				name,
				getBreakerConfig(sc),
				//d.Config().With("breaker").With("client").With(name),
				//d.Kit.G(),
			),
		),
	)

	d.httpClientMap.Store(sName, client)
	return client
}

func (d *Daenerys) HTTPServer() httpserver.Server {
	limitMaps := map[string]int{}
	for _, c := range d.config.Server.Ratelimit {
		limitMaps[c.Resource] = c.Limit
	}

	return httpserver.NewServer(
		httpserver.Name(d.localAppServiceName),
		httpserver.Port(d.config.Server.Port),
		httpserver.Tracer(d.Tracer),
		httpserver.Logger(d.Kit),
		httpserver.Tags(getServiceTags(d.config.Server.Tags)),
		httpserver.Manager(d.Manager),
		httpserver.Registry(registry.Default),
		httpserver.Ratelimit(
			ratelimit.NewConfig(
				d.Namespace,
				"",
				limitMaps,
			),
		),
		httpserver.Breaker(
			breaker.NewConfig(
				//d.Config().With("breaker").With("server"),
				//d.Kit.G(),
				d.Namespace,
				"",
				getBreakerConfigServer(d.config),
			),
		),
	)
}

func (d *Daenerys) RedisClient(service string) *redis.Redis {
	if client, ok := d.redisClientMap[service]; ok {
		return client
	}
	panic(fmt.Sprintf("namespace %s redis client for %s not exist", d.Namespace, service))
}

func (d *Daenerys) SQLClient(name string) *sql.Group {
	if client, ok := d.sqlClientMap[name]; ok {
		return client
	}
	panic(fmt.Sprintf("namespace %s sql client for %s not exist", d.Namespace, name))
}

func (d *Daenerys) KafkaConsumeClient(consumeFrom string) *kafka.KafkaConsumeClient {
	if client, ok := d.consumeClientMap[consumeFrom]; ok {
		return client
	}
	panic(fmt.Errorf("namespace %s kafka consume %s not init", d.Namespace, consumeFrom))
}

func (d *Daenerys) KafkaProducerClient(producerTo string) *kafka.KafkaClient {
	if client, ok := d.producerClientMap[producerTo]; ok {
		return client
	}
	panic(fmt.Errorf("namespace %s kafka producer %s to not init", d.Namespace, producerTo))
}

func (d *Daenerys) SyncProducerClient(producerTo string) *kafka.KafkaSyncClient {
	if client, ok := d.syncProducerClientMap[producerTo]; ok {
		return client
	}
	panic(fmt.Errorf("namespace %s kafka producer %s not init", d.Namespace, producerTo))
}

func (d *Daenerys) InitKafkaProducer(kpcList []kafka.KafkaProductConfig) error {
	for _, item := range kpcList {
		if item.UseSync {
			if _, ok := d.syncProducerClientMap[item.ProducerTo]; ok {
				continue
			}
			client, err := kafka.NewSyncProducterClient(item)
			if err != nil {
				return err
			}
			d.syncProducerClientMap[item.ProducerTo] = client
		} else {
			if _, ok := d.producerClientMap[item.ProducerTo]; ok {
				continue
			}
			client, err := kafka.NewKafkaClient(item)
			if err != nil {
				return err
			}
			d.producerClientMap[item.ProducerTo] = client
		}
	}
	return nil
}

func (d *Daenerys) InitKafkaConsume(kccList []kafka.KafkaConsumeConfig) error {
	for _, item := range kccList {
		if _, ok := d.consumeClientMap[item.ConsumeFrom]; ok {
			continue
		}
		client, err := kafka.NewKafkaConsumeClient(item)
		if err != nil {
			return err
		}
		d.consumeClientMap[item.ConsumeFrom] = client
	}
	return nil
}

func (d *Daenerys) InitRedisClient(rcList []redis.RedisConfig) error {
	for _, c := range rcList {
		cc := c
		client, err := redis.NewRedis(&cc)
		if err != nil {
			return err
		}
		d.redisClientMap[cc.ServerName] = client
	}
	return nil
}

func (d *Daenerys) InitSqlClient(sqlList []sql.SQLGroupConfig) error {
	for _, c := range sqlList {
		g, err := sql.NewGroup(c)
		if err != nil {
			return err
		}
		sql.SQLGroupManager.Add(c.Name, g)
		//err = sql.SQLGroupManager.Add(c.Name, g)
		//if err != nil {
		//	return err
		//}
		d.sqlClientMap[c.Name] = g
	}
	return nil
}

func (d *Daenerys) AddSqlClient(name string, client *sql.Group) error {
	d.sqlClientMap[name] = client
	return nil
}
func (d *Daenerys) AddRedisClient(name string, client *redis.Redis) error {
	d.redisClientMap[name] = client
	return nil
}
func (d *Daenerys) AddSyncKafkaClient(name string, client *kafka.KafkaSyncClient) error {
	d.syncProducerClientMap[name] = client
	return nil
}
func (d *Daenerys) AddAsyncKafkaClient(name string, client *kafka.KafkaClient) error {
	d.producerClientMap[name] = client
	return nil
}
func (d *Daenerys) AddHTTPClient(name string, client httpclient.Client) error {
	d.httpClientMap.Store(name, client)
	d.serverClientMap.Store(name, ServerClient{ServiceName: name})
	return nil
}
