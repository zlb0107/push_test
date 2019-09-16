package snapshot

import (
	"runtime"

	logs "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"go_common_lib/data_type"
	"go_common_lib/snapshot_pb"
	// "score/load_models"
	"strings"
	// "time"
)

// PutFeaturesToChannel 将模型特征放入通道
func PutFeaturesToChannel(version string, info data_type.LiveInfo) {
	if len(ChannelAfterVersion) >= Len_channel {
		logs.Error("channel len:", len(Channel), " too large to put")
		return
	}
	// 控制协程数量，当协程数量过多时，不再存入快照
	if runtime.NumGoroutine() > 5500 {
		logs.Error("len(ChannelAfterVersion):", len(ChannelAfterVersion), " throw out this snapshot")
		return
	}
	if len(info.OnlineFeatures) == 0 {
		return
	}

	//三部分内容：
	//1 token
	//3 online
	//5 version
	var snap proto_hall_live.SnapshotMessage
	snap.Token = info.Token

	var model proto_hall_live.ModelFeaturesInfo
	model.Type = "default"
	for idx, _ := range info.OnlineFeatures {
		model.Online = append(model.Online, &info.OnlineFeatures[idx].FeatureInfo)
	}

	snap.Snapshot = append(snap.Snapshot, &model)
	snap.Version = version

	if snap_tmp, err := proto.Marshal(&snap); err != nil {
		logs.Error("ChannelAfterVersion marshal failed and throw out this snapshot, err:", err)
	} else {
		ChannelAfterVersion <- string(snap_tmp)
	}
}

func PutPBSnapshotToChannel(info data_type.LiveInfo, version, ctrVersion, cvrVersion string) {
	if len(ChannelPB) >= Len_channel {
		logs.Error("channel len:", len(Channel), " too large to put")
		return
	}
	// 控制协程数量，当协程数量过多时，不再存入快照
	if runtime.NumGoroutine() > 5500 {
		logs.Error("len(ChannelPB):", len(ChannelPB), " throw out this snapshot")
		return
	}
	/*
		type SnapshotMessage struct {
			// 快照版本
			Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
			// token
			Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
			// 多模型
			Snapshot []*ModelFeaturesInfo `protobuf:"bytes,3,rep,name=snapshot,proto3" json:"snapshot,omitempty"`
			// 出口处最终score
			Score                float32  `protobuf:"fixed32,4,opt,name=score,proto3" json:"score,omitempty"`

	*/
	info.PBSnapShot.Version = version
	info.PBSnapShot.Token = info.Token
	info.PBSnapShot.Score = float32(info.Score)
	if len(info.PBSnapShot.Snapshot) < 2 {
		return
	}
	info.PBSnapShot.Snapshot[0].Score = float32(info.CtrScore)
	info.PBSnapShot.Snapshot[0].Type = ctrVersion
	info.PBSnapShot.Snapshot[1].Score = float32(info.CvrScore)
	info.PBSnapShot.Snapshot[1].Type = cvrVersion
	if snap_tmp, err := proto.Marshal(&info.PBSnapShot); err != nil {
		logs.Error("ChannelAfterVersion marshal failed and throw out this snapshot, err:", err)
	} else {
		ChannelPB <- string(snap_tmp)
	}
}

func appendVerifySnapshot(liveInfo *data_type.LiveInfo, ctrOrCvr, version string, feature_len int) {
	values := liveInfo.CtrNewFeatures
	idx := 0
	if ctrOrCvr == "cvr" {
		values = liveInfo.CvrNewFeatures
		idx = 1
	}
	var ymd int32
	if liveInfo.PBSnapShot.Snapshot[idx].Online == nil {
		liveInfo.PBSnapShot.Snapshot[idx].Online = make([]*proto_hall_live.FeaturesInfo, 0)
		//if len(liveInfo.PBSnapShot.Snapshot[idx].features) != 0 {
		//	ymd = liveInfo.PBSnapShot.Snapshot[idx].features[0].Ymd
		//}
	} else {
		ymd = liveInfo.PBSnapShot.Snapshot[idx].Online[0].Ymd
	}
	var mkFeatureInfo proto_hall_live.FeaturesInfo
	mkFeatureInfo.Key = "debug"

	versions := strings.Split(version, "_")
	if len(versions) < 4 {
		logs.Error("version size is < 4:", version)
		return
	}

	if versions[3] == "v5" { //走5.0之前的模型
		size := 0
		mkFeatureInfo.Indices = make([]int32, 0, size)
		mkFeatureInfo.Values = make([]float32, 0, size)
		for i, value := range values {
			if value != float64(0) {
				size += 1
				mkFeatureInfo.Indices = append(mkFeatureInfo.Indices, int32(i))
				mkFeatureInfo.Values = append(mkFeatureInfo.Values, float32(value))
			}
		}
		//} else if versions[3] == "v6" { // 走6.0的模型 valus中存的是idx:val:idx:val....
	} else {
		allSize := len(values)
		if allSize%2 != 0 {
			logs.Debug("model6.0's feature size is not even!")
			return
		}

		size := allSize / 2
		mkFeatureInfo.Indices = make([]int32, size, size)
		mkFeatureInfo.Values = make([]float32, size, size)
		for i := 0; i < size; i++ {
			mkFeatureInfo.Indices[i] = int32(values[i*2])
			mkFeatureInfo.Values[i] = (float32)(values[i*2+1])
		}
	}
	mkFeatureInfo.Size_ = int32(feature_len)
	/*
		size := len(values)
		mkFeatureInfo.Size_ = int32(size)
		mkFeatureInfo.Indices = make([]int32, size, size)
		mkFeatureInfo.Values = make([]float32, size, size)
		var i int32 = 0
		for ; i < mkFeatureInfo.Size_; i++ {
			mkFeatureInfo.Indices[i] = i
			mkFeatureInfo.Values[i] = (float32)(values[i])
		}
	*/
	mkFeatureInfo.Ymd = ymd

	liveInfo.PBSnapShot.Snapshot[idx].Online = append(liveInfo.PBSnapShot.Snapshot[idx].Online, &mkFeatureInfo)
}

// /*
// 	比PutPBSnapshotToChannel 增加了没点3:4/5/6时，打印校验快照日志
// */
// func PutPBAndVerifySnapshotToChannel(info data_type.LiveInfo, version string, ctrFeatureConf, cvrFeatureConf *load_models.FeatureConf, rec_tab string) {
// 	if len(ChannelPB) >= Len_channel {
// 		logs.Error("channel len:", len(Channel), " too large to put. rec_tab:", rec_tab, " token:", info.Token, " ctrVersion:", ctrFeatureConf.SnapshotVersion)
// 		return
// 	}
// 	// 控制协程数量，当协程数量过多时，不再存入快照
// 	if runtime.NumGoroutine() > 5500 {
// 		logs.Error("len(ChannelPB):", len(ChannelPB), " throw out this snapshot.rec_tab:", rec_tab, " token:", info.Token, " ctrVersion:", ctrFeatureConf.SnapshotVersion)
// 		return
// 	}
// 	/*
// 		type SnapshotMessage struct {
// 			// 快照版本
// 			Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
// 			// token
// 			Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
// 			// 多模型
// 			Snapshot []*ModelFeaturesInfo `protobuf:"bytes,3,rep,name=snapshot,proto3" json:"snapshot,omitempty"`
// 			// 出口处最终score
// 			Score                float32  `protobuf:"fixed32,4,opt,name=score,proto3" json:"score,omitempty"`
//
// 	*/
// 	ctrVersion := ctrFeatureConf.SnapshotVersion
// 	cvrVersion := cvrFeatureConf.SnapshotVersion
//
// 	info.PBSnapShot.Version = version
// 	info.PBSnapShot.Token = info.Token
// 	info.PBSnapShot.Score = float32(info.Score)
// 	if len(info.PBSnapShot.Snapshot) < 2 {
// 		logs.Debug("len(info.PBSnapShot.Snapshot) < 2. rec_tab:", rec_tab, " token:", info.Token, " ctrVersion:", ctrVersion)
// 		return
// 	}
// 	info.PBSnapShot.Snapshot[0].Score = float32(info.CtrScore)
// 	info.PBSnapShot.Snapshot[0].Type = ctrVersion
// 	info.PBSnapShot.Snapshot[1].Score = float32(info.CvrScore)
// 	info.PBSnapShot.Snapshot[1].Type = cvrVersion
//
// 	hour := time.Now().Hour()
// 	minute := time.Now().Minute()
// 	//3:04分打会日志  书萌用于验证数据
// 	if (rec_tab == "anonymous" || rec_tab == "big_r" || rec_tab == "register" || rec_tab == "thunder") && (hour == 3) && (minute == 4 || minute == 5 || minute == 6) {
// 		appendVerifySnapshot(&info, "ctr", ctrVersion, ctrFeatureConf.FeatureLen)
// 		appendVerifySnapshot(&info, "cvr", cvrVersion, cvrFeatureConf.FeatureLen)
// 	}
//
// 	if snap_tmp, err := proto.Marshal(&info.PBSnapShot); err != nil {
// 		logs.Error("ChannelAfterVersion marshal failed and throw out this snapshot, err:", err, " rec_tab:", rec_tab, " token:", info.Token, " ctrVersion:", ctrVersion)
// 	} else {
// 		ChannelPB <- string(snap_tmp)
// 		/*
// 			pBSnapShot := info.PBSnapShot
// 			logs.Error("info----snapshot--*******:,", info.Uid, " version:", pBSnapShot.Version, " token:", pBSnapShot.Token, " score:", pBSnapShot.Score, " ", snap_tmp)
// 			for idx, modelFeatureInfo := range pBSnapShot.Snapshot {
// 				logs.Error(" idx:", idx, "  type:", modelFeatureInfo.Type, " features:", modelFeatureInfo.Features, " features-len:", len(modelFeatureInfo.Features), "  online:", modelFeatureInfo.Online, " online-len:", len(modelFeatureInfo.Online), " trigger:", modelFeatureInfo.Trigger, " score:", modelFeatureInfo.Score)
// 			}
// 		*/
// 	}
// }
