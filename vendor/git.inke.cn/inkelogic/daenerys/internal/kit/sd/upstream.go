package sd

import (
	"errors"
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"git.inke.cn/tpc/inf/go-upstream/upstream"
	"golang.org/x/net/context"
	"strings"
)

type ikconsul struct {
	f       Factory
	cluster *upstream.Cluster
}

func Upstream(factory Factory, cluster *upstream.Cluster) core.Plugin {
	return ikconsul{
		cluster: cluster,
		f:       factory,
	}
}

const (
	resultSuccess upstream.Result = 0
	//resultTimeout                        = 1
	//resultTemporary                      = 2
	resultDNS = 3
	//resultRefused                        = 4
	resultConnectTimeout = 5
	//resultSysCall                        = 6
	resultEOF = 7
	//resultConnectUnknown                 = 99
	resultDeadline = 103
	//resultRPCUser                        = 1000
	//resultRPCUnkown                      = 1001
	resultCallError = 199
	//resultUnknown                        = 200
)

func getResultFromError(err error) upstream.Result {
	switch err := err.(type) {
	case nil:
		return resultSuccess
	default:
		if strings.Contains(err.Error(), "lookup") {
			return resultDNS
		}
		if strings.Contains(err.Error(), "canceld") {
			return resultSuccess
		}
		if strings.Contains(err.Error(), "exceed") {
			return resultDeadline
		}
		if strings.Contains(err.Error(), "timeout") {
			return resultConnectTimeout
		}
		if strings.Contains(err.Error(), "EOF") {
			return resultEOF
		}
		switch err {
		case context.DeadlineExceeded:
			return resultDeadline
		case context.Canceled:
			return resultSuccess
		default:
			return resultCallError
		}
	}
}

// TODO nil pointer
func (up ikconsul) Do(ctx context.Context, c core.Core) {
	host := up.cluster.Balancer().ChooseHost(ctx)
	if host == nil {
		c.AbortErr(errors.New("no living upstream"))
		return
	}
	defer func() {
		host.GetDetectorMonitor().PutResult(getResultFromError(c.Err()))
	}()
	plugin, err := up.f.Factory(host.Address())
	if err != nil {
		c.AbortErr(err)
		return
	}
	plugin.Do(ctx, c)
}
