// 公司内部基础库-包含了各种client以及服务发现&trace等功能
//
// Introduce
//
// 此基础库包含了redis、http、kafka的client，这些client全部集成了
// stat监控和trace系统
//
// Configure
//
// 基础库使用统一的toml格式的文件进行初始化:
//
//  [server]
//  service_name="push.backend.online"
//  port = 6100
//  	[server.http]
//   	 location="/push/send,/push/register"
//   	 logResponse="true,true"
//  [monitor]
//   alive_interval=10
//  [log]
//   level="debug"
//   rotate="dayly"
//   logpath = "logs"
//  #redis client的配置，配置项的说明请见：
//  #http://godoc.inkept.cn/git.inke.cn/BackendPlatform/golang/redis#RedisConfig
//  [[redis]]
//	 server_name="inke-redis"
//	 addr="localhost:7379"
//	 password="password"
//	 max_idle=100
//	 max_active=100
//	 idle_timeout=1000
//	 connect_timeout=1000
//	 read_timeout=1000
//	 write_timeout=1000
//	 database=0
//
package rpc // import "git.inke.cn/inkelogic/rpc-go"
