package gid

import (
	"time"
	"strconv"
	"math/rand"
	"sync"
)

var (
	g                *GId
	hostnameHashCode uint32
)

func init() {
	g = NewGId()
	hostnameHashCode = HostnameHashCode()
}

type GId struct {
	StartTime      int64
	DurationBitNum uint32
	IpBitNum       uint32
	RandBitNum     uint32
	pool           sync.Pool
}

func NewGId() *GId {
	seedGenerator := NewRand(time.Now().UnixNano())
	pool := sync.Pool{
		New: func() interface{} {
			return rand.NewSource(seedGenerator.Int63())
		},
	}
	return &GId{
		StartTime:      1546272000, //从2019-01-01 00:00:00开始计算
		DurationBitNum: 28,         //时间间隔占用28bit
		IpBitNum:       16,         //IP占用16bit
		RandBitNum:     20,         //随机数占用20bit
		pool:           pool,
	}
}

func (g *GId) Rand() uint64 {
	number := g.randNumber()
	duration := uint64(time.Now().Unix() - g.StartTime)
	return duration<<g.DurationBitNum | number>>g.DurationBitNum
}

func (g *GId) randNumber() uint64 {
	generator := g.pool.Get().(rand.Source)
	number := uint64(generator.Int63())
	g.pool.Put(generator)
	return number
}
func (g *GId) NewV1() uint64 {
	number := g.randNumber()
	duration := uint64(time.Now().Unix() - g.StartTime)
	ipCode := uint64(hostnameHashCode % (1 << g.IpBitNum))
	result := duration<<(g.IpBitNum+g.RandBitNum) | ipCode<<g.RandBitNum | (number & 0xfffff)
	return result
}
func New() uint64 {
	return g.NewV1()
}
func (g *GId) UnixFromStr(s string) int64 {
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0
	}
	return int64(v>>(g.IpBitNum+g.RandBitNum)) + g.StartTime
}
func UnixFromStr(s string) int64 {
	return g.UnixFromStr(s)

}
func (g *GId) UnixFromUint64(v uint64) int64 {
	return int64(v>>(g.IpBitNum+g.RandBitNum)) + g.StartTime
}
func UnixFromUint64(v uint64) int64 {
	return g.UnixFromUint64(v)
}
func (g *GId) StrToUint64(s string) uint64 {
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0
	}
	return v
}
func StrToUint64(s string) uint64 {
	return g.StrToUint64(s)
}
func (g *GId) FnvCodeFromStr(s string) uint32 {
	v := g.StrToUint64(s)
	return uint32((v << g.DurationBitNum) >> (g.DurationBitNum + g.RandBitNum))
}
func FnvCodeFromStr(s string) uint32 {
	return g.FnvCodeFromStr(s)
}
func (g *GId) FnvCodeFromUint64(v uint64) uint32 {
	return uint32((v << g.DurationBitNum) >> (g.DurationBitNum + g.RandBitNum))
}
func FnvCodeFromUint64(v uint64) uint32 {
	return g.FnvCodeFromUint64(v)
}
