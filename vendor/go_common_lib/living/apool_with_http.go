package living

import (
	"net/http"
	"time"
)

var Http_client *http.Client

func init() {
	t := &http.Transport{}
	//t.IdleConnTimeout = 30 * time.Second
	t.MaxIdleConnsPerHost = 10
	//t.MaxIdleConns = 20
	t.DisableKeepAlives = false
	Http_client = &http.Client{
		Transport: t,
		Timeout:   500 * time.Millisecond,
	}
}
