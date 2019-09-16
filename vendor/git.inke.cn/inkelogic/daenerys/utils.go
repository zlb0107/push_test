package daenerys

import (
	"git.inke.cn/inkelogic/daenerys/breaker"
	"path"
	"strings"
)

func getRegistryKVPath(name string) string {
	namespace := strings.Split(name, ".")[0]
	return path.Join("/service_config", namespace, name)
}

func getBreakerConfig(sc ServerClient) map[string]*breaker.BreakerConfig {
	m := make(map[string]*breaker.BreakerConfig)
	for _, b := range sc.Breaker {
		m[b.Resource] = &breaker.BreakerConfig{
			ErrorPercentThreshold:     b.EThreshold,
			ConsecutiveErrorThreshold: b.CThreshold,
			MinSamples:                b.MinSamples,
		}
	}
	return m
}

func getBreakerConfigServer(c daenerysConfig) map[string]*breaker.BreakerConfig {
	m := make(map[string]*breaker.BreakerConfig)
	for _, b := range c.Server.Breaker {
		m[b.Resource] = &breaker.BreakerConfig{
			ErrorPercentThreshold:     b.EThreshold,
			ConsecutiveErrorThreshold: b.CThreshold,
			MinSamples:                b.MinSamples,
		}
	}
	return m
}
