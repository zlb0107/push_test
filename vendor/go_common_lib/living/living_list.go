package living

import (
	logs "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"go_common_lib/go-json"
)

type LivingControl struct {
	uid_liveid_map   map[string]string
	uid_num_map      map[string]int
	uid_liveinfo_map map[string]LivingInfo
	liveid_uid_map   map[string]string
	landscap_map     map[string]int

	ObserversLock sync.RWMutex
	Observers     []interface{}
}

var Living_handler LivingControl

/*
	必须有变量ModuleName
*/
type ObserverBase interface {
	Update(this *LivingControl) error
	GetModuleName() string
}

func (this *LivingControl) AddObserver(observer interface{}) {
	this.ObserversLock.Lock()
	defer this.ObserversLock.Unlock()
	this.Observers = append(this.Observers, observer)
}

func (this *LivingControl) notify() {
	for _, observerI := range this.Observers {
		if observer, has := observerI.(ObserverBase); has {
			if observer.Update(this) != nil {
				logs.Error("update this observer failed, ModuleName:", observer.GetModuleName())
			}
		} else {
			logs.Error("notify this observer failed, ModuleName:", observer.GetModuleName())
		}
	}
}

func init() {
	go update_living()
}
func update_living() {
	Living_handler.update()
	for {
		time.Sleep(1 * time.Second)
		Living_handler.update()
	}
}

type LivingInfo struct {
	Live_id         string
	Uid             string
	NickName        string `json:"nick_name"`
	Appearance      string
	LiveTags        []*Tag `json:"live_tags,omitempty"`
	Online_num      int
	Gender          string
	Live_type       string
	SubLiveType     string `json:"sub_live_type"`
	Level           string
	Portrait        string
	LevelInt        int
	LinkMikeNum     int     `json:"link_mike_num"`
	IsBrushVideo    int     `json:"is_brush_video"`
	Quality         float64 `json:"quality"`
	ThirdPortrait   string  `json:"third_portrait"`
	IsReligion      bool    `json:"is_religion"`
	IsRetailLive    bool    `json:"is_retail_live"`    // 是否是电商直播
	PortraitHasFace bool    `json:"portrait_has_face"` // 头像是否有人像
	Landscape       int     `json:"landscape"`         //是否横屏
	RealOnlineNum   int     `json:"real_online_num"`   //真实在线人数
	WhiteListLevel  string  `json:"white_list_level"`  //默认值:-1, 不是白名单:0，白名单:1
	LongF           float64 `json:"longitude"`         //经度
	LatF            float64 `json:"latitude"`          //纬度
}
type Tag struct {
	TagId   string `json:"tag_id"`
	TagName string `json:"tag_name"`
}

func (this *LivingControl) GetAllLiveInfo() map[string]LivingInfo {
	return this.uid_liveinfo_map
}

/*
获取主播信息
*/
func (this *LivingControl) GetLiveInfo(uid string) (info LivingInfo) {
	info, _ = this.uid_liveinfo_map[uid]

	return
}

/*
	统计刷子视频的占比
*/
func (this *LivingControl) GetBrushRatio() (brushRatio float64) {
	total := len(this.uid_liveinfo_map)
	brushSize := 0
	for _, info := range this.uid_liveinfo_map {
		if (info.Live_type != "audiolive") && (info.IsBrushVideo < 0) {
			brushSize++
		}
	}
	brushRatio = float64(brushSize) / float64(total)

	return
}

/*
是否是推荐首页白名单
*/
func (this *LivingControl) IsWhiteListLevel(uid string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return false
	}

	if (info.WhiteListLevel != "0") && (info.WhiteListLevel != "-1") {
		return true
	}

	return false
}

func (this *LivingControl) GetWhiteListLevel(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return ""
	}
	return info.WhiteListLevel
}

/*
Desc：
是否是电台，包括多人电台和个人电台,常驻交友直播间(电台可抱麦新九人直播间)
找不到该主播，默认是电台
*/
func (this *LivingControl) IsAudiolive(uid string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return true
	}

	if info.Live_type == "audiolive" || info.Live_type == "residentlive" {
		return true
	}

	return false
}

/*
-2:没找到这个在线主播
-1:苍穹服务取redis中没有取到主播的实际在线人数
其他：改主播的实际在线人数
*/
func (this *LivingControl) GetRealOnlineNum(uid string) int {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return -2
	}

	return info.RealOnlineNum
}

func (this *LivingControl) GetLandscape(uid string) int {
	landscape, is_in := this.landscap_map[uid]
	if is_in == false {
		return 0
	}
	return landscape
}
func (this *LivingControl) GetAllUid() []string {
	tempList := make([]string, 0)
	for uid, _ := range this.uid_liveid_map {
		tempList = append(tempList, uid)
	}
	return tempList
}
func (this *LivingControl) GetWhitelist() []LivingInfo {
	var tempList []LivingInfo
	for _, info := range this.uid_liveinfo_map {
		if this.IsWhiteListLevel(info.Uid) {
			tempList = append(tempList, info)
		}
	}
	return tempList
}
func (this *LivingControl) Get_liveid(uid string) string {
	liveid, is_in := this.uid_liveid_map[uid]
	if is_in == false {
		return ""
	}
	return liveid
}
func (this *LivingControl) HasTag(uid, tagId string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return false
	}
	for _, tag := range info.LiveTags {
		if tag.TagId == tagId {
			return true
		}
	}
	return false
}
func (this *LivingControl) GetAllTags(uid string) []string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return nil
	}
	var tags []string
	for _, liveTag := range info.LiveTags {
		tags = append(tags, liveTag.TagName)
	}
	return tags
}
func (this *LivingControl) Get_online_num(uid string) int {
	num, is_in := this.uid_num_map[uid]
	if is_in == false {
		return -1
	}
	return num
}
func (this *LivingControl) IsBrushVideo(uid string) int {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return 0
	}
	return info.IsBrushVideo
}
func (this *LivingControl) Get_appearance(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return "-1"
	}
	return info.Appearance
}
func (this *LivingControl) Get_gender(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return "-1"
	}
	return info.Gender
}
func (this *LivingControl) GetLevel(uid string) int {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return -1
	}
	return info.LevelInt
}
func (this *LivingControl) GetPortrait(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return ""
	}
	return info.Portrait
}
func (this *LivingControl) GetSubLiveType(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in {
		return info.SubLiveType
	}
	return ""
}

// LinkMikeNum = -1表示苍穹服务没有该uid的直播信息
func (this *LivingControl) Get_live_linkMikeNum(uid string) int {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return -2
	}
	return info.LinkMikeNum
}
func (this *LivingControl) GetNickName(uid string) string {
	info, has := this.uid_liveinfo_map[uid]
	if has {
		return info.NickName
	}
	return ""
}
func (this *LivingControl) Get_live_type(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if is_in == false {
		return ""
	}
	return info.Live_type
}

func (this *LivingControl) GetQuality(uid string) float64 {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return 0.0
	}

	return info.Quality
}

func (this *LivingControl) GetThirdNonPortrait(uid string) string {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return ""
	}

	return info.ThirdPortrait
}

func (this *LivingControl) IsReligion(uid string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return false
	}

	return info.IsReligion
}

func (this *LivingControl) IsRetailLive(uid string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return false
	}

	return info.IsRetailLive
}

func (this *LivingControl) PortraitHasFace(uid string) bool {
	info, is_in := this.uid_liveinfo_map[uid]
	if !is_in {
		return false
	}

	return info.PortraitHasFace
}

func (this *LivingControl) GetUidWithLiveid(liveid string) string {
	return this.liveid_uid_map[liveid]
}

func (this *LivingControl) GetLiveTypeWithLiveid(liveid string) string {
	uid, is_in := this.liveid_uid_map[liveid]
	if !is_in {
		return ""
	}

	return this.Get_live_type(uid)
}

func (this *LivingControl) GetAppearanceWithLiveid(liveid string) string {
	uid, is_in := this.liveid_uid_map[liveid]
	if !is_in {
		return "-1"
	}

	return this.Get_appearance(uid)
}

func (this *LivingControl) GetGenderWithLiveid(liveid string) string {
	uid, is_in := this.liveid_uid_map[liveid]
	if !is_in {
		return "-1"
	}

	return this.Get_gender(uid)
}

func (this *LivingControl) GetLevelWithLiveid(liveid string) int {
	uid, is_in := this.liveid_uid_map[liveid]
	if !is_in {
		return -1
	}

	return this.GetLevel(uid)
}

func (this *LivingControl) GetPortraitWithLiveid(liveid string) string {
	uid, is_in := this.liveid_uid_map[liveid]
	if !is_in {
		return ""
	}

	return this.GetPortrait(uid)
}

func (this *LivingControl) update() {
	now := time.Now()
	resp, err := Http_client.Get("http://10.111.94.191:18089/firmament")
	if err == nil {
		defer func() { io.Copy(ioutil.Discard, resp.Body); resp.Body.Close() }()
	} else {
		logs.Error("error:", err)
		return
	}
	if resp.StatusCode != 200 {
		logs.Error("res code not 200:", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("error:", err)
		return
	}
	var lives []LivingInfo
	body_str := string(body)
	if err := json.Unmarshal([]byte(body_str), &(lives)); err != nil {
		logs.Error("Unmarshal: ", err.Error())
		return
	}
	for idx, live := range lives {
		lives[idx].LevelInt, _ = strconv.Atoi(live.Level)
	}

	uid_liveid_map := make(map[string]string)
	uid_num_map := make(map[string]int)
	uid_liveinfo_map := make(map[string]LivingInfo)
	liveid_uid_map := make(map[string]string)
	landscap_map := make(map[string]int)
	for _, live := range lives {
		uid_liveid_map[live.Uid] = live.Live_id
		uid_num_map[live.Uid] = live.Online_num
		uid_liveinfo_map[live.Uid] = live
		liveid_uid_map[live.Live_id] = live.Uid
		landscap_map[live.Uid] = live.Landscape
	}
	if len(uid_liveid_map) < 100 {
		logs.Error("all living is less 100:", len(uid_liveid_map))
		return
	}
	this.uid_liveid_map = uid_liveid_map
	this.uid_num_map = uid_num_map
	this.uid_liveinfo_map = uid_liveinfo_map
	this.liveid_uid_map = liveid_uid_map
	this.landscap_map = landscap_map
	dis := time.Since(now).Nanoseconds()
	msecond := dis / 1000000
	logs.Error("len:", len(lives), " time:", msecond, " ms")
	if dis%100 == 3 {
		//		logs.Error("len:", len(lives), " time:", msecond, " ms")
	}

	this.notify()

	logs.Flush()
}
