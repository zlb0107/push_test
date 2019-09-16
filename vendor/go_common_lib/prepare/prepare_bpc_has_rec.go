// 服务端曝光数据
package prepare

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_common_lib/data_type"
	"go_common_lib/mytime"

	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

func init() {
	Prepare_map["PrepareBpcHasRec"] = PrepareBpcHasRec{}
	logs.Info("in PrepareBpcHasRec init")
}

type PrepareBpcHasRec struct {
}

func (PrepareBpcHasRec) Get_data(request *data_type.Request, ch chan string) int {
	defer func() { ch <- "PrepareBpcHasRec" }()
	defer common.Timer("PrepareBpcHasRec", &(request.Timer_log), time.Now())

	if request == nil || request.Rec_tab == "" || request.Uid == "" || request.Session_id == "" {
		return -1
	}

	key := fmt.Sprintf("bpc_hall_%s_%s_%s", request.Rec_tab, request.Uid, request.Session_id)

	rc := BpcHasRecRedis.Get()
	defer rc.Close()

	strs, err := redis.String(rc.Do("get", key))
	if err != nil {
		return -1
	}

	// 数据格式: page_idx|uid;uid;uid
	parts := strings.Split(strs, "|")
	if len(parts) != 2 {
		return -1
	}

	pageIdx, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1
	}

	if pageIdx+1 != request.Page_idx {
		return -1
	}

	list := make(map[string]bool)
	terms := strings.Split(parts[1], ";")
	for _, term := range terms {
		if term != "" {
			list[term] = true
		}
	}

	request.BpcHasRecList = list
	return len(request.BpcHasRecList)
}
