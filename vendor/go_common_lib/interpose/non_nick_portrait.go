package interpose

import (
	//	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	logs "github.com/cihub/seelog"
	"go_common_lib/data_type"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
	"go_common_lib/mytime"
)

type NonNickPortrait struct {
	uidMap sync.Map
}
type uidStatus struct {
	NonNick     bool
	NonPortrait bool
	Timestamp   time.Time
}

// 数据59秒时就会触发刷新，1分钟时触发删除
func init() {
	var rp NonNickPortrait
	Interpose_map["NonNickPortrait"] = &rp
	logs.Warn("in NonNickPortrait init")
	go func() {
		// 每分钟一次
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rp.uidMap.Range(func(k, v interface{}) bool {
					uid := k.(string)
					uidStatus := v.(*uidStatus)
					if time.Now().Sub(uidStatus.Timestamp) > time.Minute {
						rp.uidMap.Delete(uid)
					}
					return true
				})
			}
		}
	}()
}

// 为了批量缓存uidStatus并不打乱feedid顺序，分三步走：1. 检出未缓存的uid；2. 批量缓存这些uid；3. 过滤
func (rp *NonNickPortrait) Run_interpose(request *data_type.Request) int {
	defer common.Timer("NonNickPortrait", &(request.Timer_log), time.Now())
	filter_num := len(request.Livelist)
	defer func() {
		new_log := *(request.Num_log) + " NonNickPortraitf_num:" + strconv.Itoa(filter_num)
		request.Num_log = &new_log
	}()

	noCacheUidList := make([]string, 0)
	leftList := make([]data_type.LiveInfo, 0)
	for _, info := range request.Livelist {
		if v, has := rp.uidMap.Load(info.Uid); has {
			uidStatus := v.(*uidStatus)
			// 59秒的数据就让其过期，避免下面第二遍拿时数据过期
			if time.Now().Sub(uidStatus.Timestamp) > 59*time.Second {
				noCacheUidList = append(noCacheUidList, info.Uid)
			}
		} else {
			noCacheUidList = append(noCacheUidList, info.Uid)
		}
	}

	// 缓存uid信息
	limit := 100
	length := len(noCacheUidList)
	start := 0
	end := 0
	routineNum := (length + limit - 1) / limit
	var wg sync.WaitGroup
	for i := 0; i < routineNum; i++ {
		start = i * limit
		end = start + limit
		if end > length {
			end = length
		}

		//start , end 左闭右开
		wg.Add(1)
		go func() {
			defer wg.Done()
			rp.CacheUidsInfo(noCacheUidList, start, end)
		}()
	}
	wg.Wait()

	// 过滤默认头像和默认昵称
	for _, info := range request.Livelist {
		if tmp, has := rp.uidMap.Load(info.Uid); has {
			uidStatus := tmp.(*uidStatus)
			if uidStatus.NonNick || uidStatus.NonPortrait {
				filter_num += 1
				continue
			}
			leftList = append(leftList, info)
		}
	}

	request.Livelist = leftList
	filter_num = filter_num - len(request.Livelist)
	return 0
}

func (rp *NonNickPortrait) CacheUidsInfo(noCacheUidList []string, start, end int) int {
	var uids []string
	for i := start; i < len(noCacheUidList) && i < end; i++ {
		uids = append(uids, noCacheUidList[i])
	}

	url := "http://10.111.6.202:8089/user/infos?&id="
	url += strings.Join(uids, ",")
	//resp, err := http_client_pool.Http_client.Get(url)
	body, err := http_client_pool.Get_n_url(url, 30)
	if err != nil {
		logs.Error("error:", err)
		return -1
	}
	var usersInfo usersInfoResp
	if err := json.Unmarshal(body, &usersInfo); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return -1
	}

	for _, user := range usersInfo.Users {
		var uidStatus uidStatus
		uid := strconv.FormatInt(user.Id, 10)
		if user.Nick == "Inke"+uid || user.Nick == "inke"+uid {
			uidStatus.NonNick = true
		}

		if strings.HasSuffix(user.Portrait, "MTUyODQyMzA0NTk2NiM2MjcjanBn.jpg") || user.Portrait == "" {
			uidStatus.NonPortrait = true
		}
		uidStatus.Timestamp = time.Now()
		rp.uidMap.Store(uid, &uidStatus)
	}

	return 0
}

type usersInfoResp struct {
	DmError  int     `json:"dm_error"`
	ErrorMsg string  `json:"error_msg"`
	Users    []*user `json:"users"`
}
type user struct {
	Id       int64  `json:"id"`
	Nick     string `json:"nick"`
	Portrait string `json:"portrait"`
}
