package daenerys

import (
	log "git.inke.cn/inkelogic/daenerys/log"
	"github.com/opentracing/opentracing-go"
)

type Option func(*Daenerys)

// Mode TODO
type Mode int

const (
	// TODO maybe use string better?
	Development Mode = iota // 0
	Production              // 1
)

func (m *Mode) String() string {
	switch *m {
	case Development:
		return "development"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}

func Kit(kit log.Kit) Option {
	return func(o *Daenerys) {
		o.Kit = kit
	}
}

func RunMode(mode Mode) Option {
	return func(o *Daenerys) {
		o.RunMode = mode
	}
}

func Namespace(namespace string) Option {
	return func(o *Daenerys) {
		o.Namespace = namespace
	}
}

func Sync(s bool) Option {
	return func(o *Daenerys) {
		o.Sync = s
	}
}

func Name(name string) Option {
	return func(o *Daenerys) {
		o.Name = name
	}
}

func App(app string) Option {
	return func(o *Daenerys) {
		o.App = app
	}
}

func Version(ver string) Option {
	return func(o *Daenerys) {
		o.Version = ver
	}
}

func Deps(deps string) Option {
	return func(o *Daenerys) {
		o.Deps = deps
	}
}

func Tracer(tracer opentracing.Tracer) Option {
	return func(o *Daenerys) {
		o.Tracer = tracer
	}
}

func ConfigPath(path string) Option {
	return func(o *Daenerys) {
		o.ConfigPath = path
	}
}

func ConsulAddr(addr string) Option {
	return func(o *Daenerys) {
		o.ConsulAddr = addr
	}
}

func TraceReportAddr(addr string) Option {
	return func(o *Daenerys) {
		o.TraceReportAddr = addr
	}
}
