package reader

import (
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/encoder/ini"
	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"git.inke.cn/inkelogic/daenerys/config/encoder/xml"
	"git.inke.cn/inkelogic/daenerys/config/encoder/yaml"
)

var Encoding map[string]encoder.Encoder

//默认支持的格式
func init() {
	Encoding = map[string]encoder.Encoder{
		"json": json.NewEncoder(),
		"yaml": yaml.NewEncoder(),
		"toml": toml.NewEncoder(),
		"xml":  xml.NewEncoder(),
		"yml":  yaml.NewEncoder(),
		"ini":  ini.NewEncoder(),
	}
}

//新增
func WithEncoder(e encoder.Encoder) {
	Encoding[e.String()] = e
}
