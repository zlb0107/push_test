package third

import (
	"git.inke.cn/inkelogic/daenerys/internal/core"
	"sync"
)

type Third struct {
	g *globalStage
	r *requestStage
	w *workDoneStage
}

func New() *Third {
	return &Third{
		g: &globalStage{
			ps: make([]core.Plugin, 0),
		},

		r: &requestStage{
			ps: make([]core.Plugin, 0),
		},

		w: &workDoneStage{
			ps: make([]core.Plugin, 0),
		},
	}
}

func (t *Third) OnGlobalStage() Middleware {
	return t.g
}

func (t *Third) OnRequestStage() Middleware {
	return t.r
}

func (t *Third) OnWorkDoneStage() Middleware {
	return t.w
}

type Middleware interface {
	Register([]core.Plugin)
	Stream() []core.Plugin
}

type globalStage struct {
	mu sync.Mutex
	ps []core.Plugin
}

func (g *globalStage) Register(ps []core.Plugin) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.ps = append(g.ps, ps...)
}

func (g *globalStage) Stream() []core.Plugin {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.ps
}

type requestStage struct {
	mu sync.Mutex
	ps []core.Plugin
}

func (r *requestStage) Register(ps []core.Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ps = append(r.ps, ps...)
}

func (r *requestStage) Stream() []core.Plugin {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ps
}

type workDoneStage struct {
	mu sync.Mutex
	ps []core.Plugin
}

func (w *workDoneStage) Register(ps []core.Plugin) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ps = append(w.ps, ps...)
}

func (w *workDoneStage) Stream() []core.Plugin {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.ps
}
