package discover

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	logs "github.com/cihub/seelog"

	"go_common_lib/go-json"
	"go_common_lib/http_client_pool"
)

const (
	dumpPath = "./data"
)

type DataCenterInfo struct {
	MachineCenter map[string]MachineInfo
	ServerCenter  map[string]SrvInfo
}

type MachineInfo struct {
}

type SrvInfo struct {
	IPCenter map[string]IPInfo
}

type IPInfo struct {
	IP         string
	RTT        int
	Weight     int
	RealWeight int
	CPUIdle    float64
	Expiration int64
}

func init() {
	dumpAllServices()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			dumpAllServices()
		}
	}()
}

func dumpAllServices() {
	url := "http://10.111.95.183:18097/discover/get_all"
	body, err := http_client_pool.Get_url(url)
	if err != nil {
		logs.Error("discover::dumpAllServices http_client_pool.Get_url error: ", err)
		return
	}

	dc := DataCenterInfo{}
	if err := json.Unmarshal(body, &dc); err != nil {
		logs.Error("discover::dumpAllServices json.Unmarshal error: ", err)
		return
	}

	for sname, sInfo := range dc.ServerCenter {
		var c ServiceName
		c.Name = sname
		var e IpInfo

		for ip := range sInfo.IPCenter {
			e.Ip = ip
			c.IpList = append(c.IpList, e)
		}

		if err := dumpServerResult(&c); err != nil {
			logs.Error("discover::dumpAllServices dumpServerResult error: ", err)
			continue
		}
	}
}

func dumpServerResult(c *ServiceName) error {
	if len(c.IpList) == 0 {
		return nil
	}

	_, err := os.Stat(dumpPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err = os.Mkdir(dumpPath, 0755); err != nil {
			return err
		}
	}

	tmpFile := filepath.Join(dumpPath, fmt.Sprintf("%s.tmp", c.Name))
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	desFile := filepath.Join(dumpPath, fmt.Sprintf("%s.endpoints", c.Name))
	if err := os.Rename(tmpFile, desFile); err != nil {
		return err
	}

	return nil
}

func loadServerResult(name string) (*ServiceName, error) {
	data, err := ioutil.ReadFile(filepath.Join(dumpPath, fmt.Sprintf("%s.endpoints", name)))
	if err != nil {
		return nil, err
	}

	c := &ServiceName{}
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}

type LocalDump struct {
}

var Local = &LocalDump{}

func (d *LocalDump) GetHost(sname string) (string, error) {
	c, err := loadServerResult(sname)
	if err != nil {
		return "", err
	}

	l := int64(len(c.IpList))
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(l)

	return c.IpList[idx].Ip, nil
}
