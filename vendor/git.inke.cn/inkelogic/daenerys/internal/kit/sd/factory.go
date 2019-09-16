package sd

import (
	"git.inke.cn/inkelogic/daenerys/internal/core"
)

type Factory interface {
	Factory(host string) (core.Plugin, error)
}
