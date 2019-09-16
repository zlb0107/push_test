package data_type

import (
	"context"
	"go_common_lib/debug"
	"go_common_lib/snapshot_pb"
	"strconv"
	"time"

	"go_common_lib/options"
)

const Show_info_prefix string = "show_info_"
const Show_user_prefix string = "show_user_"
const Old_user_timestamp string = "old_user_timestamp"
const Old_user_bitmap_before5 string = "old_user_bitmap_before5"
const Old_user_bitmap_after5 string = "old_user_bitmap_after5"
const Old_user_prefix string = "old_user_"

//thunder weight里面用
type TriggerScore struct {
	Key_prefix string
	Score      string
}

//该结构目前只有谛听/thunder/zeus 服务在用
type LiveInfo struct {
	Tag_id                            int            `json:"tag_id"`
	LiveId                            string         `json:"live_id"`
	Gender                            string         `json:"gender"`
	Token                             string         `json:"token"`
	Uid                               string         `json:"uid"`
	User_id                           string         `json:"User_id,omitempty"`
	Distance                          float32        `json:"distance"`
	Distance64                        float64        `json:"distance64"`
	Id                                string         `json:"id,omitempty"`
	RecTab                            string         `json:"rec_tab,omitempty"`
	Trigger                           string         `json:"-"`
	TriggerScores                     []TriggerScore `json:"-"`
	Online                            string         `json:"-"`
	Feature_redis                     string         `json:"-"`
	Feature_user_als                  string         `json:"-"`
	Feature_live_als                  string         `json:"-"`
	Score                             float64        `json:"score"`
	CtrScore                          float64
	CvrScore                          float64
	Feature                           []string   `json:"-"`
	Live_feature                      []string   `json:"-"`
	Live_feature_snapshot             string     `json:"-"`
	User_live_feature                 []string   `json:"-"`
	User_live_feature_snapshot        string     `json:"-"`
	Live_user_feature                 []string   `json:"-"`
	Live_user_feature_snapshot        string     `json:"-"`
	Realtime_feature                  []string   `json:"-"`
	Dot_feature                       []string   `json:"-"`
	Dot_feature_snapshot              string     `json:"-"`
	Trigger_feature                   [30]string `json:"-"`
	Trigger_feature_snapshot          string     `json:"-"`
	Cates                             []string   `json:"-"`
	Append_token                      string
	Tags                              []string  `json:"-"`
	Google_tag                        string    `json:"-"`
	Class                             string    `json:"-"`
	Appearance                        string    `json:"appearance,omitempty"`
	RecReason                         string    `json:"-"`
	Other_para                        OtherPara `json:"other_para"`
	NickName                          string    `json:"nick_name"`
	Portrait                          string    `json:"portrait"`
	Third_portrait                    string    `json:"third_portrait,omitempty"`
	Face                              string    `json:"face,omitempty"`
	In_black_live                     string    `json:"in_black_live_list,omitempty"`
	In_white_live                     string    `json:"in_white_live_list,omitempty"`
	Quality                           float64   `json:"quality,omitempty"`
	Live_dot_feature                  []float64 `json:"-"`
	Live_dot_feature_snapshot         string    `json:"-"`
	Live_dot_cluster_feature          []float64 `json:"-"`
	Live_dot_cluster_feature_snapshot string    `json:"-"`
	OnlineFeatures                    []Feature `json:"-"`
	NewFeatures                       []float64 `json:"-"`
	CtrNewFeatures                    []float64 `json:"-"`
	CvrNewFeatures                    []float64 `json:"-"`
	CtrDnnFeatures                    []float64 `json:"-"`
	CvrDnnFeatures                    []float64 `json:"-"`
	IsPk                              bool
	OfflineSnapshot                   string                          `json:"-"`
	OnlineSnapshot                    string                          `json:"-"`
	NextSnapshotOffline               string                          `json:"-"`
	NextSnapshotOnline                string                          `json:"-"`
	PlanId                            string                          `json:"-"`
	PBSnapShot                        proto_hall_live.SnapshotMessage `json:"-"`
	HasWatermark                      bool                            `json:"has_watermark,omitempty"`
	FeedType                          int                             `json:"feed_type"`
	AtlasHasFilter                    string                          `json:"-"`
	Pos                               string                          `json:"-"`
	CtrLeafNewFeatures                []float64                       `json:"-"`
	CvrLeafNewFeatures                []float64                       `json:"-"`
	IsLiving                          bool                            `json:"is_living,omitempty"`   //是否开播
	DynamicCoverDistance              float32                         `json:"dc_distance,omitempty"` //动态封面的距离

	LongF float64 `json:"-"` //经度
	LatF  float64 `json:"-"` //纬度
	IsNew string  //是否是新主播
}
type OtherPara struct {
	Reason    string `json:"reason"`
	Multicard []Card `json:"multicard,omitempty"`
}
type Card struct {
	Tabkey string `json:"tabkey"`
	Tag    string `json:"tag"`
	Uid    string `json:"uid"`
	Liveid string `json:"liveid"`
}
type Feature struct {
	FeatureInfo proto_hall_live.FeaturesInfo
	Dim         string
}
type Response struct {
	Rec_tab    string     `json:"rec_tab"`
	Livelist   []LiveInfo `json:"list"`
	Error_msg  string     `json:"error_msg"`
	Error_code int        `json:"error_code"`
}
type Request struct {
	Timer_log                          *string `json:"-"`
	Num_log                            *string `json:"-"`
	Uid                                string
	Session_id                         string `json:"-"` // 这次请求的会话ID
	Scope                              int    //每一位表示一个小时内，如第一位表示1小时内，第二位表示1-2小时内，7二进制111，表示3小时内所有,特殊：0表示所有，默认1
	Infolist                           []BeatHeartInfo
	Gender                             string
	BigRLevel                          string `json:"-"`
	Tab_key                            string
	Logid                              string  `json:"-"`
	Mylogid                            string  `json:"-"`
	Longitude                          string  `json:"-"`
	Latitude                           string  `json:"-"`
	LongF                              float64 `json:"-"` //经度
	LatF                               float64 `json:"-"` //纬度
	EventTime                          string  `json:"-"`
	Livelist                           []LiveInfo
	Count                              int `json:"-"`
	Real_count                         int `json:"-"`
	Raw_len                            int
	PlanNum                            int
	UserLevel                          int
	CardPos                            int
	Score                              int
	Cates                              []string
	Expids                             []string
	Ids                                *[]string
	Less_love_uids                     map[string]float64
	Has_in_uids                        map[string]bool
	Has_shown_uids                     map[string]bool
	HasOutIds                          map[string]bool
	User_blacklist                     map[string]bool
	Live_blacklist                     map[string]bool
	Strong_filter_uids                 map[string]bool
	Has_rec_list                       map[string]bool
	BpcHasRecList                      map[string]bool
	Page0HasRecList                    map[string]bool
	Dislike_follow_list                map[string]bool
	Dislike_hall_list                  map[string]bool
	Follow_list                        map[string]bool
	Film_has_shown_list                map[string]bool
	Is_new                             bool
	IsReturn30                         bool
	IsNotActive                        bool
	IsCut                              bool `json:"-"`
	IsBigR                             bool
	IsBlacklist                        bool
	IsShuazi                           bool
	IsMetis                            bool           // 公演服务请求
	Cross_class                        string         `json:"-"`
	Class                              string         `json:"-"`
	Filter_map                         map[string]int `json:"-"`
	User_feature                       []string       `json:"-"`
	User_feature_snapshot              string         `json:"-"`
	User_hot_cluster                   string         `json:"-"`
	Paras                              *string        `json:"-"`
	Rec_tab                            string         `json:"-"`
	Normal_result                      string         `json:"-"`
	Hades_result                       string         `json:"-"`
	Bname                              string         `json:"-"`
	Zid_str                            string         `json:"-"`
	Types                              map[string]int `json:"-"`
	Page_idx                           int            `json:"-"`
	User_dot_features_snapshot         string         `json:"-"`
	User_dot_features                  []float64      `json:"-"`
	User_dot_cluster_features_snapshot string         `json:"-"`
	User_cluster_dot_features          []float64      `json:"-"`
	StartTime                          string         `json:"-"`
	EndTime                            string         `json:"-"`
	PlanID                             string         `json:"-"`
	User                               string         `json:"-"`
	Event                              string         `json:"-"`
	Amount                             string         `json:"-"`
	Pos                                int            `json:"-"`
	//PaoPaoHistoryChan                  chan PaoPaoHistory
	PaoPaoHistorys      map[string]bool
	PaoPaoTodayRedisKey string
	PaoPaoTodayRecNum   int
	City                string `json:"-"`
	Prov                string `json:"-"`
	SnapshotVersion     string `json:"-"`
	DebugCtx            *debug.DebugContext
	Ctx                 context.Context `json:"-"`
	Interest_sex        string          `json:"-"`
	Is_rookie           string          `json:"-"`

	TraceId      string                        `json:"-"`
	BetaParams   map[string]options.BetaParams `json:"-"`
	FeedbackList []LiveInfo                    `json:"-"`
	SupLivelist  []LiveInfo                    `json:"-"`
	AttrMap      map[string]interface{}        `json:"-"`
	OtherAttrMap map[string]interface{}
}

type Info struct {
	Obj_id  string
	Obj_uid string
	Pos     string
	Token   string
	Obj_key string
}
type MdEinfo struct {
	Infos   []Info
	Tab_key string
}
type BeatHeartInfo struct {
	Uid         string
	Md_eid      string
	Record_time string
	Md_mod      string
	Men_time    time.Time
	Md_einfo    MdEinfo
	Raw_str     string `json:"-"`
	Logid       string
	Vv          string
	Cv          string
}
type PosType []Info

func (c PosType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c PosType) Len() int {
	return len(c)
}
func (c PosType) Less(i, j int) bool {
	pi, _ := strconv.Atoi(c[i].Pos)
	pj, _ := strconv.Atoi(c[j].Pos)
	return pi < pj
}

type SortType []BeatHeartInfo

func (c SortType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c SortType) Len() int {
	return len(c)
}
func (c SortType) Less(i, j int) bool {
	pi, _ := strconv.ParseInt(c[i].Record_time, 10, 0)
	pj, _ := strconv.ParseInt(c[j].Record_time, 10, 0)
	return pi > pj
}

type ChanStruct struct {
	Livelist []LiveInfo
	Name     string
	Rate     float64
}

const Prefer_name string = "prefer"
const Red_packets_name string = "red_packets"
const Rec_name string = "personal"
const Hot_name string = "hot"
const New_name string = "new"
const Old_name string = "old"
const Talent_name string = "talent"
