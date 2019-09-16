package store

import (
	logs "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"go_common_lib/data_type"
	"go_common_lib/feed_store/pb"
	"strconv"
)

var reqChan chan *data_type.Request

func init() {
	reqChan = make(chan *data_type.Request, 5000)
	go func() {
		for {
			req := <-reqChan
			storeShown(req)
		}
	}()
}
func GetKey(uid string) string {
	return "shown_" + uid
}
func StoreShown(req *data_type.Request) {
	reqChan <- req
}

const QUEUELEN int = 3000

func storeShown(req *data_type.Request) {
	deleteNum := 0
	totalNum := len(req.Livelist) + len(*(req.Ids))
	if totalNum > 3000 {
		deleteNum = totalNum - 3000
	}
	key := GetKey(req.Uid)
	storeInfo := shown.Shown{}
	storeInfo.Ids = make([]uint64, 0)
	for idx, idStr := range *(req.Ids) {
		if idx < deleteNum {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			logs.Error("err:", err)
			continue
		}
		storeInfo.Ids = append(storeInfo.Ids, id)
	}
	for _, info := range req.Livelist {
		id, err := strconv.ParseUint(info.LiveId, 10, 64)
		if err != nil {
			logs.Error("err:", err)
			continue
		}
		storeInfo.Ids = append(storeInfo.Ids, id)
	}
	value, err := proto.Marshal(&storeInfo)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	rc := StoreRedis.Get()
	defer rc.Close()
	rc.Do("setex", key, 86400*5, value)
}
func GetShown(uid string) *[]string {
	key := GetKey(uid)
	rc := StoreRedis.Get()
	defer rc.Close()
	showList := make([]string, 0)
	value, err := rc.Do("get", key)
	if err != nil {
		logs.Error("err:", err)
		return &showList
	}
	if value == nil {
		//logs.Error("value is nil")
		return &showList
	}
	storeInfo := shown.Shown{}
	err = proto.Unmarshal([]byte(value.([]byte)), &storeInfo)
	if err != nil {
		logs.Error("err:", err)
		return &showList
	}
	for _, id := range storeInfo.Ids {
		showList = append(showList, strconv.FormatUint(id, 10))
	}
	return &showList
}
