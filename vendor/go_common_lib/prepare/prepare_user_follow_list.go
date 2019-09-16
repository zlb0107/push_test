package prepare

import (
	logs "github.com/cihub/seelog"
	"go_common_lib/connect_pool"
	"go_common_lib/connect_pool/gen-go/following"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"strconv"
	"time"
)

type FollowList struct {
}

func init() {
	var rp FollowList
	Prepare_map["FollowList"] = rp
	logs.Warn("in FollowList init")
}

func (rp FollowList) Get_data(request *data_type.Request, c chan string) int {
	defer func() { c <- "FollowList" }()
	is_hit := false
	defer common.TimerV2("FollowList", &(request.Timer_log), time.Now(), &is_hit)
	pooledClient, err := connect_pool.Follow_con_pool.Get()
	if err != nil {
		logs.Error("pooled Client failed:", err)
		return -1
	}
	defer pooledClient.Close()
	client := pooledClient.RawClient().(*following.RelationServiceClient)
	uid, trans_err := strconv.ParseInt(request.Uid, 10, 64)
	if trans_err != nil {
		logs.Error("trans_err:", trans_err)
		return -1
	}
	res, get_err := client.GetAllFollowings(uid)
	if get_err != nil {
		pooledClient.MarkUnusable()
		logs.Error("Error get_err: ", get_err)
		return -1
	}
	request.Follow_list = make(map[string]bool)
	for _, uid := range res.GetUids() {
		uid_str := strconv.FormatInt(uid, 10)
		request.Follow_list[uid_str] = true
	}

	return 0
}
