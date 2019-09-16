package discover

import (
	"time"

	"go_common_lib/data_type"
	"go_common_lib/http_client_pool"
)

type Discover interface {
	GetHost(name string) (string, error)
}

type IpResult struct {
	Ip  string
	Err error
}

func GetUrl(sname string, timeout int) (string, error) {
	var d Discover
	if sname == "score" {
		d = CPUIdle
	} else {
		d = RTT
	}

	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	ch := make(chan IpResult, 1)

	go func() {
		var r IpResult
		defer func() { ch <- r }()

		ip, err := d.GetHost(sname)
		r.Ip = ip
		r.Err = err
	}()

	select {
	case ret := <-ch:
		if ret.Err == nil {
			return ret.Ip, nil
		}

	case <-timer.C:
	}

	d = Local
	return d.GetHost(sname)
}

func GetResult(sname string, req *data_type.Request, urlPostfix string, timeout int) ([]byte, error) {
	ip, err := GetUrl(sname, 10)
	if err != nil {
		return nil, err
	}

	url := "http://" + ip + urlPostfix
	//记录时间起点
	start := time.Now()
	var useTime int64

	defer func() {
		//记录所用时间,ms
		TimeChan <- SnameTimeInfo{sname: sname, time: useTime, ip: ip}
	}()

	body, err := http_client_pool.Get_n_url(url, timeout)
	if err != nil {
		useTime = 1000
	} else {
		useTime = time.Since(start).Nanoseconds() / 1000000
	}

	return body, err
}
