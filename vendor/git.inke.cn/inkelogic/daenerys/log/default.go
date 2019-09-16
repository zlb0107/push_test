package log

import (
	"bytes"
	"fmt"
	"git.inke.cn/BackendPlatform/golang/logging"
	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
	"strings"
)

type defaultLogger struct {
	log *logging.Logger
}

func Default(path string) *defaultLogger {
	d := logging.New()
	d.SetFlags(0)
	d.SetPrintLevel(false)
	d.SetHighlighting(false)
	d.SetRotateByHour()
	d.SetOutputByName(path)
	return &defaultLogger{d}
}

func (d *defaultLogger) Rotate(rotate string) *defaultLogger {
	if rotate == "day" {
		d.log.SetRotateByDay()
		return d
	}
	return d
}

func (d *defaultLogger) Level(level string) *defaultLogger {
	d.log.SetLevelByString(level)
	return d
}

func (d *defaultLogger) Logger() Logger {
	return With(defaultLogger{d.log})
}

func (d defaultLogger) Log(kvs ...interface{}) error {
	var buf = &bytes.Buffer{}
	logMap := make(map[string]interface{}, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		k := fmt.Sprintf("%v", kvs[i])
		value := kvs[i+1]
		if v, ok := value.(error); ok {
			logMap[k] = v.Error()
		} else {
			logMap[k] = value
		}
	}
	b, _ := json.NewEncoder().Encode(logMap)
	buf.Write(b)
	d.log.Info(strings.TrimSpace(buf.String()))
	buf.Reset()
	return nil
}
