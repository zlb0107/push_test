package feature_prepare

import (
	"runtime"
	"strings"
	"time"

	logs "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"go_common_lib/data_type"
	"go_common_lib/mytime"
	"go_common_lib/snapshot_pb"
)

type OnlineFeaturePrepare struct {
}

func init() {
	var rp OnlineFeaturePrepare
	PrepareMap["OnlineFeaturePrepare"] = rp
	logs.Warn("in OnlineFeaturePrepare init")
}
func (rp OnlineFeaturePrepare) GetData(req *data_type.Request, c chan FeatureWrapperStruct) int {
	var featureWrapper FeatureWrapperStruct
	featureWrapper.Name = ONLINEFEATURE
	defer func() { c <- featureWrapper }()

	// 黑名单不取在线特征打快照
	if req.IsBlacklist {
		return -1
	}

	if runtime.NumGoroutine() > 5000 {
		logs.Error("runtime.NumGoroutine() > 5000, OnlineFeaturePrepare stop GetData.")
		return -1
	}

	var ch = make(chan struct{}, 1)
	var retCode int
	go func() {
		retCode = getData(req, &featureWrapper)
		ch <- struct{}{}
	}()
	select {
	case <-time.After(50 * time.Millisecond):
		logs.Error("OnlineFeaturePrepare timeout.")
		return -1
	case <-ch:
		return retCode
	}
}

func getData(req *data_type.Request, featureWrapper *FeatureWrapperStruct) int {
	defer common.Timer(featureWrapper.Name, &(req.Timer_log), time.Now())
	GetCommonData(req, "online", featureWrapper)
	//feature_prepare.ShowFeatureWrapperStruct(featureWrapper)

	for _, newFeature := range featureWrapper.NewFeatures {
		// 每个newFeature中判断是在线特征还是离线特征
		if strings.HasSuffix(newFeature.GetType, "online") {
			// 遍历一个特征中的所有redis特征
			for _, feature := range newFeature.Features {
				// 判断每一个redis特征Dim
				if feature.Dim == "user" || feature.Dim == "model" {
					// 给每个live都添加
					for idx, _ := range req.Livelist {
						if len(feature.Features) != 1 {
							logs.Error("解析快照错误")
							return -1
						}
						tmp := []byte(feature.Features[0])
						var fInfo proto_hall_live.FeaturesInfo
						err := proto.Unmarshal(tmp, &fInfo)
						if err != nil {
							logs.Error(err)
						}
						req.Livelist[idx].OnlineFeatures = append(req.Livelist[idx].OnlineFeatures, data_type.Feature{
							FeatureInfo: fInfo,
							Dim:         feature.Dim,
						})
					}
				} else {
					// 某些mget会超时，超时时跳过该redis特征，写入其他特征
					if len(req.Livelist) > len(feature.Features) {
						continue
					}
					for idx, _ := range req.Livelist {
						tmp := []byte(feature.Features[idx])
						var fInfo proto_hall_live.FeaturesInfo
						err := proto.Unmarshal(tmp, &fInfo)
						if err != nil {
							logs.Error(err)
						}

						req.Livelist[idx].OnlineFeatures = append(req.Livelist[idx].OnlineFeatures, data_type.Feature{
							FeatureInfo: fInfo,
							Dim:         feature.Dim,
						})
					}
				}
			}
		} else {
			logs.Error("Offline features do not require a snapshot.")
		}

		req.SnapshotVersion = featureWrapper.SnapshotVersion
	}

	logs.Flush()
	return 0
}
