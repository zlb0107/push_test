package redis

import (
	"fmt"

	"git.inke.cn/BackendPlatform/miniredis"
)

// NewMockRedis returns MiniRedis and close func and err.
// NOTE: mockRedis can't auto expire,you can use `TTL key`
// to replace auto expire
func NewMockRedis() (mockRedis *Redis, closeFunc func(), err error) {
	mockServer, err := miniredis.Run()
	closeFunc = func() {
		mockServer.Close()
	}
	if err != nil {
		return
	}
	mockRedis, err = NewRedis(&RedisConfig{
		Addr:        fmt.Sprintf("localhost:%s", mockServer.Port()),
		ServerName:  "mock.redis",
		ReadTimeout: 1000,
	})
	return
}
