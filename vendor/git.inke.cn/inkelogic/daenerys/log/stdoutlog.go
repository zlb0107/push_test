package log

import (
	"bytes"
	"fmt"
	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
	"os"
	"strings"
)

type stdoutLogger struct{}

func Stdout() Logger {
	return With(stdoutLogger{})
}

func (d stdoutLogger) Log(kvs ...interface{}) error {
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
	fmt.Fprint(os.Stdout, strings.TrimSpace(buf.String())+"\n")
	buf.Reset()
	return nil
}
