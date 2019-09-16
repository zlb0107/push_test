// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: features.proto

package proto_hall_live

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import encoding_binary "encoding/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FeaturesInfo struct {
	// redis key
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// 特征长度，无论稀疏特征，还是稠密特征都设置
	Size_ int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// 特征值 f1,f2,f3 ...
	Values []float32 `protobuf:"fixed32,3,rep,packed,name=values" json:"values,omitempty"`
	// 特征索引稀疏特征时设置，稠密特征默认为空
	Indices []int32 `protobuf:"varint,4,rep,packed,name=indices" json:"indices,omitempty"`
	// 特征时间，用于在指定分区解析离线特征
	Ymd                  int32    `protobuf:"varint,5,opt,name=ymd,proto3" json:"ymd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeaturesInfo) Reset()         { *m = FeaturesInfo{} }
func (m *FeaturesInfo) String() string { return proto.CompactTextString(m) }
func (*FeaturesInfo) ProtoMessage()    {}
func (*FeaturesInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_features_eab1d656adbbf00c, []int{0}
}
func (m *FeaturesInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeaturesInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeaturesInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *FeaturesInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeaturesInfo.Merge(dst, src)
}
func (m *FeaturesInfo) XXX_Size() int {
	return m.Size()
}
func (m *FeaturesInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FeaturesInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FeaturesInfo proto.InternalMessageInfo

func (m *FeaturesInfo) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *FeaturesInfo) GetSize_() int32 {
	if m != nil {
		return m.Size_
	}
	return 0
}

func (m *FeaturesInfo) GetValues() []float32 {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *FeaturesInfo) GetIndices() []int32 {
	if m != nil {
		return m.Indices
	}
	return nil
}

func (m *FeaturesInfo) GetYmd() int32 {
	if m != nil {
		return m.Ymd
	}
	return 0
}

type OfflineFeaturesInfo struct {
	// 用于存储快照数据，服务取之即用
	Snapshot []*FeaturesInfo `protobuf:"bytes,1,rep,name=snapshot" json:"snapshot,omitempty"`
	// 特征值
	Features []*FeaturesInfo `protobuf:"bytes,2,rep,name=features" json:"features,omitempty"`
	// 写入redis的key
	Key                  string   `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OfflineFeaturesInfo) Reset()         { *m = OfflineFeaturesInfo{} }
func (m *OfflineFeaturesInfo) String() string { return proto.CompactTextString(m) }
func (*OfflineFeaturesInfo) ProtoMessage()    {}
func (*OfflineFeaturesInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_features_eab1d656adbbf00c, []int{1}
}
func (m *OfflineFeaturesInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OfflineFeaturesInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OfflineFeaturesInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *OfflineFeaturesInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OfflineFeaturesInfo.Merge(dst, src)
}
func (m *OfflineFeaturesInfo) XXX_Size() int {
	return m.Size()
}
func (m *OfflineFeaturesInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_OfflineFeaturesInfo.DiscardUnknown(m)
}

var xxx_messageInfo_OfflineFeaturesInfo proto.InternalMessageInfo

func (m *OfflineFeaturesInfo) GetSnapshot() []*FeaturesInfo {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

func (m *OfflineFeaturesInfo) GetFeatures() []*FeaturesInfo {
	if m != nil {
		return m.Features
	}
	return nil
}

func (m *OfflineFeaturesInfo) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func init() {
	proto.RegisterType((*FeaturesInfo)(nil), "FeaturesInfo")
	proto.RegisterType((*OfflineFeaturesInfo)(nil), "OfflineFeaturesInfo")
}
func (m *FeaturesInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeaturesInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if m.Size_ != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(m.Size_))
	}
	if len(m.Values) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(len(m.Values)*4))
		for _, num := range m.Values {
			f1 := math.Float32bits(float32(num))
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f1))
			i += 4
		}
	}
	if len(m.Indices) > 0 {
		dAtA3 := make([]byte, len(m.Indices)*10)
		var j2 int
		for _, num1 := range m.Indices {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0x22
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	if m.Ymd != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(m.Ymd))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *OfflineFeaturesInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OfflineFeaturesInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Snapshot) > 0 {
		for _, msg := range m.Snapshot {
			dAtA[i] = 0xa
			i++
			i = encodeVarintFeatures(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Features) > 0 {
		for _, msg := range m.Features {
			dAtA[i] = 0x12
			i++
			i = encodeVarintFeatures(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Key) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintFeatures(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintFeatures(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *FeaturesInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovFeatures(uint64(l))
	}
	if m.Size_ != 0 {
		n += 1 + sovFeatures(uint64(m.Size_))
	}
	if len(m.Values) > 0 {
		n += 1 + sovFeatures(uint64(len(m.Values)*4)) + len(m.Values)*4
	}
	if len(m.Indices) > 0 {
		l = 0
		for _, e := range m.Indices {
			l += sovFeatures(uint64(e))
		}
		n += 1 + sovFeatures(uint64(l)) + l
	}
	if m.Ymd != 0 {
		n += 1 + sovFeatures(uint64(m.Ymd))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *OfflineFeaturesInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Snapshot) > 0 {
		for _, e := range m.Snapshot {
			l = e.Size()
			n += 1 + l + sovFeatures(uint64(l))
		}
	}
	if len(m.Features) > 0 {
		for _, e := range m.Features {
			l = e.Size()
			n += 1 + l + sovFeatures(uint64(l))
		}
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovFeatures(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovFeatures(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozFeatures(x uint64) (n int) {
	return sovFeatures(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FeaturesInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeatures
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FeaturesInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeaturesInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeatures
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			m.Size_ = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Size_ |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType == 5 {
				var v uint32
				if (iNdEx + 4) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
				iNdEx += 4
				v2 := float32(math.Float32frombits(v))
				m.Values = append(m.Values, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFeatures
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthFeatures
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 4
				if elementCount != 0 && len(m.Values) == 0 {
					m.Values = make([]float32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					if (iNdEx + 4) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
					iNdEx += 4
					v2 := float32(math.Float32frombits(v))
					m.Values = append(m.Values, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Values", wireType)
			}
		case 4:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFeatures
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Indices = append(m.Indices, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFeatures
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthFeatures
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Indices) == 0 {
					m.Indices = make([]int32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFeatures
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Indices = append(m.Indices, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Indices", wireType)
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ymd", wireType)
			}
			m.Ymd = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Ymd |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeatures(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFeatures
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *OfflineFeaturesInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeatures
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: OfflineFeaturesInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OfflineFeaturesInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Snapshot", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeatures
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Snapshot = append(m.Snapshot, &FeaturesInfo{})
			if err := m.Snapshot[len(m.Snapshot)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Features", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeatures
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Features = append(m.Features, &FeaturesInfo{})
			if err := m.Features[len(m.Features)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeatures
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeatures(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFeatures
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipFeatures(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFeatures
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFeatures
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthFeatures
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowFeatures
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipFeatures(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthFeatures = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFeatures   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("features.proto", fileDescriptor_features_eab1d656adbbf00c) }

var fileDescriptor_features_eab1d656adbbf00c = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x4b, 0x4d, 0x2c,
	0x29, 0x2d, 0x4a, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xaa, 0xe0, 0xe2, 0x71, 0x83,
	0x8a, 0x78, 0xe6, 0xa5, 0xe5, 0x0b, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x06, 0x81, 0x98, 0x42, 0x42, 0x5c, 0x2c, 0xc5, 0x99, 0x55, 0xa9, 0x12, 0x4c, 0x0a,
	0x8c, 0x1a, 0xac, 0x41, 0x60, 0xb6, 0x90, 0x18, 0x17, 0x5b, 0x59, 0x62, 0x4e, 0x69, 0x6a, 0xb1,
	0x04, 0xb3, 0x02, 0xb3, 0x06, 0x53, 0x10, 0x94, 0x27, 0x24, 0xc1, 0xc5, 0x9e, 0x99, 0x97, 0x92,
	0x99, 0x9c, 0x5a, 0x2c, 0xc1, 0xa2, 0xc0, 0xac, 0xc1, 0x1a, 0x04, 0xe3, 0x82, 0xcc, 0xad, 0xcc,
	0x4d, 0x91, 0x60, 0x05, 0x1b, 0x02, 0x62, 0x2a, 0xd5, 0x72, 0x09, 0xfb, 0xa7, 0xa5, 0xe5, 0x64,
	0xe6, 0xa5, 0xa2, 0x38, 0x40, 0x93, 0x8b, 0xa3, 0x38, 0x2f, 0xb1, 0xa0, 0x38, 0x23, 0xbf, 0x44,
	0x82, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x88, 0x57, 0x0f, 0x59, 0x41, 0x10, 0x5c, 0x1a, 0xa4, 0x14,
	0xe6, 0x1b, 0x09, 0x26, 0xac, 0x4a, 0x61, 0xd2, 0x30, 0x6f, 0x31, 0xc3, 0xbd, 0xe5, 0xe4, 0x7c,
	0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7,
	0xc0, 0xa5, 0x90, 0x9c, 0x9f, 0xab, 0x97, 0x99, 0x97, 0x9d, 0xaa, 0x57, 0x9c, 0x5a, 0x54, 0x96,
	0x99, 0x9c, 0x0a, 0x09, 0x22, 0xbd, 0x8c, 0xc4, 0x9c, 0x1c, 0xbd, 0x9c, 0xcc, 0xb2, 0xd4, 0x28,
	0x7e, 0x34, 0x81, 0x24, 0x36, 0xb0, 0x80, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x69, 0x1a,
	0xef, 0x56, 0x01, 0x00, 0x00,
}