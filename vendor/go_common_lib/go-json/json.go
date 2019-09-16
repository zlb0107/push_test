package json

import (
	jsoniter "github.com/json-iterator/go"
)

var JsonConfig = jsoniter.Config{
	EscapeHTML:                    true,
	SortMapKeys:                   true,
	ValidateJsonRawMessage:        true,
	ObjectFieldMustBeSimpleString: true,
}.Froze()

func Unmarshal(data []byte, v interface{}) error {
	return JsonConfig.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return JsonConfig.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return JsonConfig.MarshalToString(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return JsonConfig.MarshalIndent(v, prefix, indent)
}
