package http_client_pool

import (
	"net/http"
	"time"
)

var Http_client *http.Client
var Http_50_client *http.Client
var Http_20_client *http.Client
var ClientMap map[int]*http.Client

func initClient(n int) *http.Client {
	t := &http.Transport{}
	//t.IdleConnTimeout = 30 * time.Second
	t.MaxIdleConnsPerHost = 10
	//t.MaxIdleConns = 20
	t.DisableKeepAlives = false
	return &http.Client{
		Transport: t,
		Timeout:   time.Duration(n) * time.Millisecond,
	}
}
func init() {
	ClientMap = make(map[int]*http.Client)
	t := &http.Transport{}
	//t.IdleConnTimeout = 30 * time.Second
	t.MaxIdleConnsPerHost = 10
	//t.MaxIdleConns = 20
	t.DisableKeepAlives = false
	Http_client = &http.Client{
		Transport: t,
		Timeout:   500 * time.Millisecond,
	}
	t1 := &http.Transport{}
	t1.MaxIdleConnsPerHost = 10
	t1.DisableKeepAlives = false
	Http_50_client = &http.Client{
		Transport: t1,
		Timeout:   50 * time.Millisecond,
	}
	t2 := &http.Transport{}
	t2.MaxIdleConnsPerHost = 10
	t2.DisableKeepAlives = false
	Http_20_client = &http.Client{
		Transport: t2,
		Timeout:   20 * time.Millisecond,
	}

}
