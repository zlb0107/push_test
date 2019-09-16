package daenerys

import (
	"flag"
	"git.inke.cn/inkelogic/daenerys/config"
	"os"
)

var (
	consulAddr      string
	traceReportAddr string
)

const (
	LOG_ROTATE_HOUR  = "hour"
	LOG_ROTATE_DAY   = "day"
	LOG_ROTATE_MONTH = "month"
)

func init() {
	var (
		fallbackConsulAddr      = "127.0.0.1:8500"
		fallbackTraceReportAddr = "127.0.0.1:6831"
	)

	if addr, ok := os.LookupEnv("CONSUL_ADDR"); ok {
		fallbackConsulAddr = addr
	}
	if addr, ok := os.LookupEnv("TRACE_ADDR"); ok {
		fallbackTraceReportAddr = addr
	}
	flag.StringVar(&consulAddr, "consul-addr", fallbackConsulAddr, "consul agent http addr.")
	flag.StringVar(&traceReportAddr, "trace-addr", fallbackTraceReportAddr, "trace agent udp report addr.")
	config.ConsulAddr = consulAddr
}
