package discover

const SNAMELISTPATH = "/discover/service_name_list"
const SNAMEIPDIR = "/discover/"

type IpInfo struct {
	Ip     string
	Weight int
}

type ServiceName struct {
	Name    string
	Explain string
	IpList  []IpInfo
}

type ServiceNameList struct {
	List []ServiceName
}

type SnameTimeInfo struct {
	sname string
	time  int64
	ip    string
}

type ServiceInfo struct {
	Sname          string
	MachineMap     map[string]*Machine
	AverageNumbers []float64
}

type Machine struct {
	IP       string
	RTTs     []float64 // 统计平响数据
	CPUIdles []float64 // 统计cpu idle数据
	α        int       // 正样本数
	β        int       // 负样本数
}
