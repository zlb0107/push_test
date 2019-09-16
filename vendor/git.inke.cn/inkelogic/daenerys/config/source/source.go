// Package source is the interface for sources
package source

import (
	"crypto/md5" // #nosec
	"fmt"
	"time"
)

// 数据源提供的方法
type Source interface {
	Read() (*ChangeSet, error)
	Watch() (Watcher, error)
	String() string
}

// 数据源监视器
type Watcher interface {
	Next() (*ChangeSet, error)
	Stop() error
}

// 数据源被解析后的结构
type ChangeSet struct {
	Data      []byte
	Checksum  string
	Format    string
	Source    string
	Timestamp time.Time
}

func (c *ChangeSet) Sum() string {
	h := md5.New() // #nosec
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}