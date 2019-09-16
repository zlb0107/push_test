package upstream

import (
	log "git.inke.cn/BackendPlatform/golang/logging"
)

var (
	logging *log.Logger
)

func init() {
	logging = log.New()
}

func SetLogger(l *log.Logger) {
	logging = l
}
