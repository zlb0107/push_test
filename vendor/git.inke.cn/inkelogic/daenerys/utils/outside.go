package utils

import (
	"errors"
	"fmt"
	"git.inke.cn/BackendPlatform/golang/ecode"

	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"

	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"git.inke.cn/inkelogic/daenerys/internal/kit/retry"
	"git.inke.cn/inkelogic/daenerys/internal/kit/sd"
	"git.inke.cn/tpc/inf/go-upstream/config"
	"git.inke.cn/tpc/inf/go-upstream/registry"
)

func Register(manager *registry.ServiceManager, appServiceName string, protoType string, tags map[string]string, ip string, port int) (*config.Register, error) {
	var err error
	name := fmt.Sprintf("%s-%s", appServiceName, protoType)
	cfg := config.NewRegister(name, ip, port)
	cfg.ServiceTags = tags
	cfg.TagsWatchPath, err = sd.RegistryKVPath(appServiceName, "/service_tags")
	if err != nil {
		return nil, err
	}
	if manager != nil {
		if err := manager.Register(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func MakeAppServiceName(app, name string) string {
	if len(app) == 0 {
		return name
	}
	return app + "." + name
}

func LastError(err error) error {
	switch e := err.(type) {
	case retry.RetryError:
		err = e.Final
	}
	return err
}

func Base64(buf []byte) string {
	//udp 65535byte limit
	limit := 20480
	if len(buf) == 0 || len(buf) > limit {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf)
}
func DumpRespBody(resp *http.Response) []byte {
	if resp == nil {
		return nil
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(buf) == 0 {
		return nil
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return buf
}

var (
	URIDecoder = schema.NewDecoder()
	Valid      = validator.New()
)

func init() {
	URIDecoder.IgnoreUnknownKeys(true)
}

func Bind(raw *http.Request, model interface{}, obj ...interface{}) error {
	bodyBytes := make([]byte, 0)
	if raw.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(raw.Body)
		defer func() {
			raw.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}()
	}
	method := raw.Method
	reqUrl := raw.URL.String()
	switch method {
	case "GET":
		dErr := URIDecoder.Decode(model, raw.URL.Query())
		vErr := Valid.Struct(model)
		if dErr != nil || vErr != nil {
			//estr := fmt.Sprintf("http parse reqUrl failed:|reqUrl:%v|err:%v,%v|", reqUrl, dErr, vErr)
			//return errors.New(estr)
			return ecode.ParamErr
		}
	case "POST":
		uErr := json.NewEncoder().Decode(bodyBytes, &model)
		vErr := Valid.Struct(model)
		if uErr != nil || vErr != nil {
			estr := fmt.Sprintf("http parse body failed:|reqUrl:%v|body:%s|err:%v,%v|", reqUrl, bodyBytes, uErr, vErr)
			return errors.New(estr)
		}
	default:
	}
	if len(obj) > 0 {
		dErr := URIDecoder.Decode(obj[0], raw.URL.Query())
		vErr := Valid.Struct(obj[0])
		if dErr != nil || vErr != nil {
			estr := fmt.Sprintf("http parse atom failed:|reqUrl:%v|err:%v,%v|", reqUrl, dErr, vErr)
			return errors.New(estr)
		}
	}
	return nil
}

type WrapResp struct {
	Code int         `json:"dm_error"`
	Msg  string      `json:"error_msg"`
	Data interface{} `json:"data"`
}

func NewWrapResp(data interface{}, err error) WrapResp {
	e := ecode.Cause(err)
	return WrapResp{
		Code: e.Code(),
		Msg:  e.Message(),
		Data: data,
	}
}
