package toml

import (
	"errors"
	"git.inke.cn/inkelogic/daenerys/config/encoder"
	"git.inke.cn/inkelogic/daenerys/config/encoder/toml"
	"git.inke.cn/inkelogic/daenerys/config/reader"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"github.com/imdario/mergo"
	"time"
)

type tomlReader struct {
	tm encoder.Encoder
}

func NewReader() reader.Reader {
	return &tomlReader{
		tm: toml.NewEncoder(),
	}
}

func (t *tomlReader) Merge(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	var merged map[string]interface{}

	for _, m := range changes {
		if m == nil {
			continue
		}
		if len(m.Data) == 0 {
			continue
		}
		//选择文件的编码方式
		codec, ok := reader.Encoding[m.Format]
		if !ok {
			codec = t.tm
		}

		var data map[string]interface{}
		if err := codec.Decode(m.Data, &data); err != nil {
			return nil, err
		}
		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	b, err := t.tm.Encode(merged)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Data:      b,
		Source:    "toml",
		Format:    t.tm.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (t *tomlReader) Values(ch *source.ChangeSet) (reader.Values, error) {
	if ch == nil {
		return nil, errors.New("changeset is nil")
	}
	if ch.Format != "toml" {
		return nil, errors.New("unsupported format")
	}
	return newValues(ch)
}

func (t *tomlReader) String() string {
	return "toml"
}
