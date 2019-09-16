package redis_pool

import (
	"time"

	"go_common_lib/config"

	logs "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type RedisInfo struct {
	Instance      **redis.Pool
	Host          string
	Auth          string
	Maxidle       int
	Maxactive     int
	Idle_timeout  int
	Read_timeout  time.Duration
	Write_timeout time.Duration
}

type RedisInfo_v2 struct {
	Instance      **redis.Pool
	Host          string
	Auth          string
	Maxidle       int
	Maxactive     int
	Idle_timeout  int
	Conn_timeout  time.Duration
	Read_timeout  time.Duration
	Write_timeout time.Duration
}

func RedisInit(ri RedisInfo) {
	// 建立连接池
	*(ri.Instance) = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     config.AppConfig.DefaultInt("redis.maxidle", ri.Maxidle),
		MaxActive:   config.AppConfig.DefaultInt("redis.maxactive", ri.Maxactive),
		IdleTimeout: time.Duration(ri.Idle_timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ri.Host, redis.DialReadTimeout(ri.Read_timeout*time.Nanosecond), redis.DialWriteTimeout(ri.Write_timeout*time.Nanosecond))
			if err != nil {
				logs.Error("dial failed host:", ri.Host, " auth:", ri.Auth)
				return nil, err
			}
			if _, err := c.Do("AUTH", ri.Auth); err != nil {
				c.Close()
				logs.Error("auth failed, host:", ri.Host, " auth:", ri.Auth)
				return nil, err
			}
			return c, nil
		},
	}
}

//V2 有连接超时选项
func RedisInit_v2(ri RedisInfo_v2) {
	// 建立连接池
	*(ri.Instance) = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     config.AppConfig.DefaultInt("redis.maxidle", ri.Maxidle),
		MaxActive:   config.AppConfig.DefaultInt("redis.maxactive", ri.Maxactive),
		IdleTimeout: time.Duration(ri.Idle_timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ri.Host,
				redis.DialConnectTimeout(ri.Conn_timeout*time.Nanosecond),
				redis.DialReadTimeout(ri.Read_timeout*time.Nanosecond),
				redis.DialWriteTimeout(ri.Write_timeout*time.Nanosecond))
			if err != nil {
				logs.Error("dial failed host:", ri.Host, " auth:", ri.Auth)
				return nil, err
			}
			if _, err := c.Do("AUTH", ri.Auth); err != nil {
				c.Close()
				logs.Error("auth failed host:", ri.Host, " auth:", ri.Auth)
				return nil, err
			}
			return c, nil
		},
	}
}
