package filter

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/pk"
	"strings"
)

type FilterPk struct {
}

func init() {
	var rp FilterPk
	Filter_map["FilterPk"] = rp
	logs.Warn("in FilterPk init")
}
func (rp FilterPk) Filter_live(info *data_type.LiveInfo, req *data_type.Request) bool {
	// 新开播主播扶持（只陪你聊）过滤pk
	if len(info.Token) > 13 && info.Token[:13] == "rec_7_48_1_0^" {
		return pk.Pk_handler.Is_pk(info.Uid)
	}

	//跳舞的主播需要过滤正在pk
	//多品类跳舞token：rec_7_18_3_0^373401653_1559360144388_60^93^82|10001
	var isDancing bool = false

	if len(info.Token) > 8 {
		if info.Token[:8] == "rec_7_24" {
			isDancing = true
		} else {
			tokens := strings.Split(info.Token, "^")
			if len(tokens) > 2 {
				for _, t := range tokens[2:] {
					if t == "93" {
						isDancing = true
					}
				}
			}
		}
	}

	if isDancing {
		return pk.Pk_handler.Is_pk(info.Uid)
	}
	return false
}
