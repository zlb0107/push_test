package fetchattr

import (
	"go_common_lib/data_type"
)

type AttrFetcher interface {
	Get_attr(*data_type.Request, chan string) int
}

var Fetcher_map map[string]AttrFetcher

//推荐架构可配置
var FetcherCfg_map map[string]func(moduleName string, params map[string]interface{}) AttrFetcher

func init() {
	Fetcher_map = make(map[string]AttrFetcher)
	FetcherCfg_map = make(map[string]func(moduleName string, params map[string]interface{}) AttrFetcher)
}
