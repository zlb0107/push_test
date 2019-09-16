package new_sups

import (
	"go_common_lib/data_type"
	"go_common_lib/load_confs"
)

type Supplement interface {
	Get_list(*data_type.Request, chan []data_type.LiveInfo, *load_confs.SupInfo) int
}

type SupplementCfg interface {
	Supplement
	GetSupID() int
	GetSupLimitNum() int
}

var NewSupMap map[string]Supplement
var NewSupCfgMap map[string]func(string, map[string]interface{}) SupplementCfg

func init() {
	NewSupMap = make(map[string]Supplement)
	NewSupCfgMap = make(map[string]func(string, map[string]interface{}) SupplementCfg)
}
