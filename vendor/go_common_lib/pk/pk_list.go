package pk

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	//	"strconv"
	"time"
)

const LIMIT int = 50

type PkControl struct {
	pk_map  map[string]bool
	pk_list [LIMIT]string
}

var Pk_handler PkControl

func init() {
	Pk_handler.pk_map = make(map[string]bool)
	Pk_handler.pk_list = [LIMIT]string{}
	go update_pk_list()
}
func update_pk_list() {
	Pk_handler.update()
	for {
		time.Sleep(1 * time.Second)
		Pk_handler.update()
	}
}

func (this *PkControl) GetList() [LIMIT]string {
	return this.pk_list
}
func (this *PkControl) Is_pk(uid string) bool {
	_, is_in := this.pk_map[uid]
	if is_in == false {
		return false
	}
	return true
}

func (this *PkControl) update() {
	now := time.Now()
	rc := Pk_redis.Get()
	defer rc.Close()
	values, err := redis.Strings(rc.Do("smembers", "pklist"))
	if err != nil {
		logs.Error("get redis failed:", err)
		return
	}
	temp_map := make(map[string]bool)
	temp_list := [LIMIT]string{}
	for idx, value := range values {
		temp_map[value] = true
		if idx >= LIMIT {
			continue
		}
		temp_list[idx] = value
	}
	//	//debug
	//	for i := 1; i < 10; i++ {
	//		temp_map[strconv.Itoa(i)] = true
	//	}
	this.pk_map = temp_map
	this.pk_list = temp_list
	dis := time.Since(now).Nanoseconds()
	msecond := dis / 1000000
	if dis%1000 == 3 {
		logs.Error("pk_len:", len(this.pk_map), " time:", msecond, " ms")
	}
	logs.Flush()
}
