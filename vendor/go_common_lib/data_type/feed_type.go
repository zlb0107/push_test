package data_type

import "sync"

const (
	FEEDREC    string = "feed_rec"
	FEEDNEAR   string = "feed_near"
	FEEDFOLLOW string = "feed_follow"
)

type FeedInfo struct {
	FeedID         string           `json:"EntityId"`
	Status         int              `json:"acl_type"`
	UID            string           `json:"EntityOwnerId"`
	Block          int              `json:"block"`
	Delete         int              `json:"delete"`
	Attachments    []AttachmentInfo `json:"attachments"`
	Type           int              `json:"type"`
	HasWatermark   bool             `json:"watermark,omitempty"` // 是否有水印
	VideoLevel     string           //视频清晰度,用于过滤
	AtlasHasFilter string           // “”:空代表没取到或者feedid不是图集或者是正常的图集   0:代表视频有问题，1:标签有问题 2:带声乐的图集问题,不满足条件，需要过滤  10:图集满足条件,不需要过滤;其中0，1，2是书萌那边在redis中存的黑名单中的值,兼容
	Pos            string           `json:"pos,omitempty"` // 分发评级, 没有时默认为3(即正常发布), wiki: http://wiki.inkept.cn/pages/viewpage.action?pageId=50824524
}
type FeedKafkaInfo struct {
	FeedID    uint64   `json:"EntityId"`
	UID       uint64   `json:"EntityOwnerId"`
	Ext       ExtInfo  `json:"Ext"`
	EventType string   `json:"EventType"`
	Data      DataInfo `json:"Data"`
	Atom      AtomInfo `json:"Atom"`
}
type DataInfo struct {
	Type    int         `json:"type"`
	Content ContentInfo `json:"content"`
}
type AtomInfo struct {
	Pos   string `json:"pos"`
	Token string `json:"token"`
}
type ContentInfo struct {
	Attachments []AttachmentInfo `json:"attachments"`
}
type AttachmentInfo struct {
	Type int           `json:"type"`
	Data InterDataInfo `json:"data"`
}
type InterDataInfo struct {
	Duration int `json:"duration"`
}
type SimilarFeedKafkaInfo struct {
	FeedID string      `json:"entity_id"`
	Ext    SimiExtInfo `json:"ext"`
}
type SimiExtInfo struct {
	RepeatList string `json:"repeat_list"`
}
type ExtInfo struct {
	IsRepeat   string   `json:"is_repeat"`
	RepeatList []string `json:"repeat_list"`
	Status     int      `json:"acl_type"`
	Block      int      `json:"block"`
}
type GroupInfo struct {
	FeedsMap sync.Map
	Leader   string
	Status   int
	Block    int
}
type FeedResponse struct {
	Rec_tab    string             `json:"rec_tab"`
	List       []FeedResponseInfo `json:"list"`
	Error_msg  string             `json:"error_msg"`
	Error_code int                `json:"error_code"`
}
type FeedResponseInfo struct {
	Pos        int       `json:"pos"`
	FeedId     string    `json:"feed_id"`
	Uid        string    `json:"uid"`
	Token      string    `json:"token"`
	Distance   float32   `json:"distance"`
	Reason     string    `json:"reason"`
	Other_para OtherPara `json:"other_para"`
}

//phoenix 使用
type PhoenixResp struct {
	LeaderId       string `json:"leaderid"`
	Uid            string `json:"uid"`
	FeedType       int    `json:"feed_type"`
	HasWatermark   bool   `json:"has_watermark"`
	AtlasHasFilter string `json:"atlas_has_filter"` // “”:空代表没取到或者feedid不是图集或者是正常的图集   0:代表视频有问题，1:标签有问题 2:带声乐的图集问题,不满足条件，需要过滤  10:图集满足条件,不需要过滤;其中0，1，2是书萌那边在redis中存的黑名单中的值,兼容
	Pos            string `json:"pos,omitempty"`
}
