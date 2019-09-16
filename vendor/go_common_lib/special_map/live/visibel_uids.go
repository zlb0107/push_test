package live_special_map

import (
	"bytes"
	"context"
	"git.inke.cn/inkelogic/rpc-go"
	logs "github.com/cihub/seelog"
	"go_common_lib/living"
	"go_common_lib/special_map"
	"strconv"
	"time"

	"go_common_lib/go-json"
)

type NotVisibleUids struct {
	special_map.SpecialUids
}

var NotVisibleUidsController NotVisibleUids

type HttpInfos struct {
	Visible_infos []VisibleInfo `json:"infos"`
}
type VisibleInfo struct {
	Visible int `json:"is_visible"`
	Active  int `json:"is_active"`
	Uid     int `json:"uid"`
}

func init() {
	go func() {
		for {
			//获得全量在线直播
			allUids := living.Living_handler.GetAllUid()
			NotVisibleUidsController.UpdateHttp(allUids, 500, "/near/visible_infos?uid=", get_visible)
			time.Sleep(5 * time.Second)
		}
	}()
}
func get_visible(notVisibleMap *map[string]bool, uri string) {
	serviceName := "user.location.user-near"
	result, err := rpc.HttpPost(context.TODO(), serviceName, uri, nil, bytes.NewReader([]byte(`{"offset":0, "type":153}`)))
	if err != nil {
		logs.Error("post fail:" + err.Error() + string(result))
		return
	}

	var stat_json HttpInfos
	if err := json.Unmarshal(result, &stat_json); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return
	}
	for _, info := range stat_json.Visible_infos {
		if info.Visible != 0 || info.Active != 0 {
			continue
		}
		(*notVisibleMap)[strconv.Itoa(info.Uid)] = true
	}
}
