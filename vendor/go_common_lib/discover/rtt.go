package discover

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"

	logs "github.com/cihub/seelog"
)

type RTTInfo struct {
	ServiceMap sync.Map
}

var RTT = &RTTInfo{}

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			RTT.Update()
		}
	}()
}

func (this *RTTInfo) Update() {
	this.ServiceMap.Range(func(k, v interface{}) bool {
		sname := k.(string)
		this.updateService(sname)
		return true
	})
}

func (this *RTTInfo) updateService(sname string) {
	v, has := this.ServiceMap.Load(sname)
	serviceInfo := ServiceInfo{}
	if has {
		serviceInfo = v.(ServiceInfo)
	}

	var snameMap map[string]float64
	machineMap := make(map[string]*Machine)

	snameMap = this.GetSingleServerRTTInfo(sname)
	if snameMap == nil {
		return
	}

	sum := 0.0
	for _, rtt := range snameMap {
		sum += rtt
	}
	average := sum / float64(len(snameMap))
	// 第一次填充历史记录
	if len(serviceInfo.AverageNumbers) == 0 {
		serviceInfo.AverageNumbers = make([]float64, CountWindos)
		for i := range serviceInfo.AverageNumbers {
			serviceInfo.AverageNumbers[i] = average
		}
	}

	for ip, rtt := range snameMap {
		m := &Machine{IP: ip}
		m0, has := serviceInfo.MachineMap[ip]
		if has {
			m.RTTs = append(m0.RTTs[1:], rtt)
		} else {
			m.RTTs = serviceInfo.AverageNumbers
		}
		for idx, rtt := range m.RTTs {
			if rtt <= serviceInfo.AverageNumbers[idx] {
				m.α += step
			} else {
				m.β += step
			}
		}

		machineMap[ip] = m
	}

	serviceInfo.MachineMap = machineMap

	// 更新服务信息
	this.ServiceMap.Store(sname, serviceInfo)
}

func (this *RTTInfo) GetSingleServerRTTInfo(sname string) map[string]float64 {
	url := "http://10.111.95.183:18097/discover/get_rtt_list?server_name=" + sname
	body, err := http_client_pool.Get_url(url)
	if err != nil {
		logs.Error("discover::RTT GetSingleServerRTTInfo http.Get error: ", err, ", sname: ", sname)
		return nil
	}

	tempMap := make(map[string]float64)
	err = json.Unmarshal(body, &tempMap)
	if err != nil {
		logs.Error("discover::RTT GetSingleServerInfo json.Unmarshal err: ", err, ", sname: ", sname)
		return nil
	}

	if len(tempMap) == 0 {
		logs.Error("discover::RTT GetSingleServerInfo len(tempMap) == 0, sname: ", sname)
		return nil
	}

	return tempMap
}

func (this *RTTInfo) GetHost(sname string) (string, error) {
	rand.Seed(time.Now().Unix())
	v, has := this.ServiceMap.Load(sname)
	if !has {
		this.updateService(sname)
		v, has = this.ServiceMap.Load(sname)
		if !has {
			return "", fmt.Errorf("not in server map, sname: %v", sname)
		}
	}

	serviceInfo := v.(ServiceInfo)
	maxScore := 0.0
	ip := ""
	for _, m := range serviceInfo.MachineMap {
		mBeta := NextBeta(float64(m.α), float64(m.β))
		if maxScore < mBeta {
			maxScore = mBeta
			ip = m.IP
		}
	}

	return ip, nil
}
