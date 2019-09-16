package breaker

import (
	"sync"
)

var circuitSetting = new(sync.Map) // string -> Setting

type BreakerConfig struct {
	ErrorPercentThreshold     int  `toml:"error_percent_threshold"`
	ConsecutiveErrorThreshold int  `toml:"connsecutive_error_threshold"`
	MinSamples                int  `toml:"minsamples"`
	Break                     bool `toml:"break"`
}

func Configure(name string, config *BreakerConfig) {
	circuitSetting.Store(name, config)
}

func getSetting(name string) *BreakerConfig {
	setting, ok := circuitSetting.Load(name)
	if !ok {
		return nil
	}
	return setting.(*BreakerConfig)
}
