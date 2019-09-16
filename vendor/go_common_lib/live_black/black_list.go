package live_black

import (
	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

//./redis-cli -h r-2zee82b066955c84359.redis.rds.aliyuncs.com -a vOByljzlvh26 getbit  black_live_bitmap 215428
type BlackMap struct {
	BlackUidMap map[string]bool
	Str_after5  string
	Str_before5 string
}

var Bm BlackMap

func init() {
	go update_bitmap()
}
func update_bitmap() {
	Bm.updateNew()
	for {
		time.Sleep(1 * time.Second)
		Bm.updateNew()
	}
}

func (this *BlackMap) is_setbit(suid string, str string) bool {
	num, err := strconv.Atoi(suid)
	if err != nil {
		logs.Error("suid:is not num", err, " suid", suid)
		return false
	}
	//取到的字符串，每8位为一组，一个字节,每个字节从后往前存位
	block := num / 8
	bit_offset := (7 - num%8)
	if len(str) <= block {
		//logs.Error("len(str):", len(str), " block:", block, " num:", num)
		return false
	}
	mybyte := str[block]
	if (mybyte & (1 << uint32(bit_offset))) != 0 {
		return true
	}
	return false
}
func (this *BlackMap) Bad(uid string) bool {
	//	if len(uid) < 5 {
	//		return false
	//	}
	//	uid_after5 := uid[5:len(uid)]
	//	uid_before5 := uid[0:5]
	//	if this.is_setbit(uid_after5, this.Str_after5) && this.is_setbit(uid_before5, this.Str_before5) {
	//		return true
	//	}
	//	return false
	return this.BlackUidMap[uid]
}
func (this *BlackMap) updateNew() {
	rc := Black_redis.Get()
	defer rc.Close()
	now := time.Now()
	timestamp, err := redis.String(rc.Do("get", "black_live_timestamp"))
	if err != nil {
		logs.Error("update black bitmap failed:", err)
		return
	}
	str, err := redis.String(rc.Do("get", "black_live_list_"+timestamp))
	if err != nil {
		logs.Error("update black bitmap failed:", err)
		return
	}
	terms := strings.Split(str, ";")
	tempMap := make(map[string]bool)
	for _, term := range terms {
		if term == "" {
			continue
		}
		tempMap[term] = true
	}
	Bm.BlackUidMap = tempMap
	dis := time.Since(now).Nanoseconds()
	msecond := dis / 1000000
	if dis%100 == 3 {
		logs.Error("len:", len(tempMap), " time:", msecond, " ms")
	}
	logs.Flush()
}
func (this *BlackMap) update() {
	rc := Black_redis.Get()
	defer rc.Close()
	now := time.Now()
	timestamp, err := redis.String(rc.Do("get", "black_live_timestamp"))
	if err != nil {
		logs.Error("update black bitmap failed:", err)
		return
	}
	str_before5, err := redis.String(rc.Do("get", "black_live_bitmap_before5"+timestamp))
	if err != nil {
		logs.Error("update black bitmap failed:", err)
		return
	}
	str_after5, err := redis.String(rc.Do("get", "black_live_bitmap_after5"+timestamp))
	if err != nil {
		logs.Error("update black bitmap failed:", err)
		return
	}
	Bm.Str_after5 = str_after5
	Bm.Str_before5 = str_before5
	dis := time.Since(now).Nanoseconds()
	msecond := dis / 1000000
	if dis%100 == 3 {
		logs.Error("len:", len(str_before5), " ", len(str_after5), " time:", msecond, " ms")
	}
	logs.Flush()
}
