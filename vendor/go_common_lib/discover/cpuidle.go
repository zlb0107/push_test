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

// CountWindos CountWindos是收集计算平响的窗口,值越大分流量层越明显,平响越收敛
const CountWindos = 1200

// step step越大,汤姆森采样对平响的变动反应越大
const step = 1

type CPUIdleInfo struct {
	ServiceMap sync.Map
}

var CPUIdle = &CPUIdleInfo{}

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			CPUIdle.Update()
		}
	}()
}

// updateServiceMap 异步更新
func (this *CPUIdleInfo) Update() {
	this.ServiceMap.Range(func(k, v interface{}) bool {
		sname := k.(string)
		this.updateService(sname)
		return true
	})
}

func (this *CPUIdleInfo) updateService(sname string) {
	v, has := this.ServiceMap.Load(sname)
	serviceInfo := ServiceInfo{}
	if has {
		serviceInfo = v.(ServiceInfo)
	}

	var snameMap map[string]float64
	machineMap := make(map[string]*Machine)

	// 使用cpu idle
	snameMap = this.GetSingleServerCPUIdleInfo(sname)
	if snameMap == nil {
		return
	}

	sum := 0.0
	for _, cpu_idle := range snameMap {
		sum += cpu_idle
	}
	average := sum / float64(len(snameMap))
	// 第一次填充历史记录
	if len(serviceInfo.AverageNumbers) == 0 {
		serviceInfo.AverageNumbers = make([]float64, CountWindos)
		for i := range serviceInfo.AverageNumbers {
			serviceInfo.AverageNumbers[i] = average
		}
	}

	for ip, cpu_idle0 := range snameMap {
		m := &Machine{IP: ip}
		m0, has := serviceInfo.MachineMap[ip]
		if has {
			m.CPUIdles = append(m0.CPUIdles[1:], cpu_idle0)
		} else {
			m.CPUIdles = serviceInfo.AverageNumbers
		}
		for idx, cpu_idle := range m.CPUIdles {
			if cpu_idle >= serviceInfo.AverageNumbers[idx] {
				m.α += step
			} else {
				m.β += step
			}
		}

		// logs.Debugf("discover更新: service_name=%s, ip=%s, α=%d, β=%d", sname, m.IP, m.α, m.β)
		machineMap[ip] = m
	}

	serviceInfo.MachineMap = machineMap

	// 更新服务信息
	this.ServiceMap.Store(sname, serviceInfo)
}

func (this *CPUIdleInfo) GetSingleServerCPUIdleInfo(sname string) map[string]float64 {
	url := "http://10.111.95.183:18097/discover/get_cpu_idle_list?server_name=" + sname
	body, err := http_client_pool.Get_url(url)
	if err != nil {
		logs.Error("discover::CPUIdle GetSingleServerCPUIdleInfo http.Get error: ", err, ", sname: ", sname)
		return nil
	}

	tempMap := make(map[string]float64)
	err = json.Unmarshal(body, &tempMap)
	if err != nil {
		logs.Error("discover::CPUIdle GetSingleServerCPUIdleInfo json.Unmarshal error: ", err, ", sname: ", sname)
		return nil
	}

	if len(tempMap) == 0 {
		logs.Error("discover::CPUIdle GetSingleServerCPUIdleInfo len(tempMap) == 0, sname: ", sname)
		return nil
	}

	return tempMap
}

func (this *CPUIdleInfo) GetHost(sname string) (string, error) {
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
