package rpc

import (
	"fmt"
	"git.inke.cn/tpc/inf/go-upstream/registry"
	"strings"
	"sync"

	"git.inke.cn/inkelogic/daenerys"
	dutil "git.inke.cn/inkelogic/daenerys/utils"
	"git.inke.cn/inkelogic/rpc-go/naming"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
)

var (
	kvWatchedMutex = new(sync.Mutex)
	kvWatched      = make(map[string]bool)
	kvmessage      = make(chan *naming.KvMessage)
)

const (
	HTTP_PROTO  string = "http"
	RPC_PROTO   string = "rpc"
	HTTPS_PROTO string = "https"
)

var (
	NOT_FOUND_SERVICE string = "balance no found service:%v, check has not init"
)

var (
	SERVICE_REMOTE_CONFIG_PRE = "/service_config/"
)

func getRegistryKVPath(path string) (string, error) {
	serviceName := GetServiceName()
	if len(serviceName) == 0 {
		return "", fmt.Errorf("consul kv path, servicename:%v;", serviceName)
	}
	namespace := strings.Split(serviceName, ".")[0]
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = fmt.Sprintf("%s%s/%s%s", SERVICE_REMOTE_CONFIG_PRE, namespace, serviceName, path)
	return path, nil
}

func GetRemoteKvNoWatch(path string) (string, error) {
	npath, err := getRegistryKVPath(path)
	if err != nil {
		return "", err
	}
	value, _, err := registry.Default.ReadManual(npath)
	return value, err
}

func GetRemoteKv(path string) (string, error) {
	watchPath, err := getRegistryKVPath(path)
	if err != nil {
		return "", err
	}
	value, _, err := registry.Default.ReadManual(watchPath)
	kvWatchedMutex.Lock()
	if _, ok := kvWatched[path]; !ok {
		go func() {
			last := value
			watchCh := registry.Default.WatchManual(watchPath)
			for data := range watchCh {
				if data != last {
					kvmessage <- &naming.KvMessage{
						Path:       path,
						OriginPath: path,
						PathValue:  data,
					}
					last = data
				}
			}
		}()
	}
	kvWatchedMutex.Unlock()
	return value, err
}

func NextKvMessage() (*naming.KvMessage, error) {
	u := <-kvmessage
	return u, nil
}

func GetServiceConfig(name string) (ServerClient, error) {
	return serviceClient(name, nil)
}

func GetCluster(name string) *upstream.Cluster {
	sc, _ := serviceClient(name, nil)
	return daenerys.Default.Clusters.Cluster(sc.Cluster.Name)
}

func RegisterService(proto string, port int, tags []string) {
	appServiceName := dutil.MakeAppServiceName(daenerys.Default.App, daenerys.Default.Name)
	_, err := dutil.Register(daenerys.Default.Manager, appServiceName, proto, getServiceTags(tags), config.LocalIPString(), port)
	if err != nil {
		panic(err)
	}
}

func serviceClient(service string, config RequestOptionIntercace) (ServerClient, error) {
	sc, err := daenerys.Default.FindServerClient(service)
	if err != nil {
		return sc, fmt.Errorf(NOT_FOUND_SERVICE, service)
	}

	if config == nil {
		return sc, nil
	}
	if config.GetTimeOut() > 0 {
		sc.ReadTimeout = config.GetTimeOut()
	}
	if config.GetRetryTimes() > 0 {
		sc.RetryTimes = config.GetRetryTimes()
	}
	if config.GetSlowTime() > 0 {
		sc.SlowTime = config.GetSlowTime()
	}
	return sc, nil
}
