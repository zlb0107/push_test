package debug

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	"go_common_lib/debug/struct"
	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
)

const (
	// ES_API ES接口
	ES_API = "http://10.111.70.183:9201"
	// ES_Index 索引
	ES_Index_FORMAT = "rec_debug_%s_%s"
	// ES_ID_FORMAT uid+时间戳
	ES_ID_FORMAT = "%s_%d"
)

// ES_Type 每个服务有不同的ES_Type用来标识数据
var ES_Type = "rec"

// DebugContext Debug上下文，会存储到ES中，请使用NewDebugContext创建Debug上下文。
type DebugContext struct {
	Uid       string                         `json:"uid"`
	Timestamp int64                          `json:"timestamp"`
	Data      []byte                         `json:"data"`
	data      proto_debugctx.DebugCtxMessage `json:"-"`
}

// Init 初始化方法，使用前必须初始化
func Init(serviceName string) {
	ES_Type = serviceName
}

// NewDebugContext 新建一个Debug上下文
func NewDebugContext(uid string) *DebugContext {
	return &DebugContext{Uid: uid, Timestamp: time.Now().Unix()}
}

func (d *DebugContext) DebugSupplement(listData *proto_debugctx.ListData) {
	d.data.SupplementDatas = append(d.data.SupplementDatas, listData)
}

func (d *DebugContext) DebugMerge(listData *proto_debugctx.ListData) {
	d.data.MergeData = listData
}

func (d *DebugContext) DebugScatter(listData *proto_debugctx.ListData) {
	d.data.ScatterDatas = append(d.data.ScatterDatas, listData)
}

func (d *DebugContext) DebugInterpose(listData *proto_debugctx.ListData) {
	d.data.InterposeDatas = append(d.data.InterposeDatas, listData)
}

func (d *DebugContext) DebugOutput(listData *proto_debugctx.ListData) {
	d.data.OutputData = listData
}

// PutES 将内存中的数据存入ES
//索引格式：rec_debug_zeus_<date>
//type: zeus
func (d *DebugContext) PutES() error {
	now := time.Now()
	index := fmt.Sprintf(ES_Index_FORMAT, ES_Type, now.Format("20060102"))
	_id := fmt.Sprintf(ES_ID_FORMAT, d.Uid, now.UnixNano()/1000000)
	//   http://xx.xx.xx.xx/rec_debug_zeus_20060102/zeus/10000_1557244800000
	url := fmt.Sprintf("%s/%s/%s/%s", ES_API, index, ES_Type, _id)
	d.Data, _ = proto.Marshal(&d.data)
	buf, err := json.Marshal(d)
	if err != nil {
		return err
	}
	_, err = http_client_pool.GetUrlPostN(url, bytes.NewReader(buf), 1000)
	if err != nil {
		return err
	}
	return nil
}
