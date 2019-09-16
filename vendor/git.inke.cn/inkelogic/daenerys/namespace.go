package daenerys

import (
	log "git.inke.cn/BackendPlatform/golang/logging"
	ns "git.inke.cn/inkelogic/daenerys/internal/kit/namespace"
	"golang.org/x/net/context"
	"os"
	"path"
	"sync"
)

var GlobalNamespace = newNamespace()

type namespace struct {
	daenerys sync.Map
}

func newNamespace() *namespace {
	return &namespace{}
}

func (n *namespace) Add(namespace string, opts ...Option) {
	if _, ok := n.daenerys.Load(namespace); ok {
		return
	}
	opts = append(opts, Namespace(namespace))
	d := New()
	d.Init(opts...)
	n.daenerys.LoadOrStore(namespace, d)
}

func (n *namespace) Get(namespace string) *Daenerys {
	d, _ := n.daenerys.Load(namespace)
	if d == nil {
		return nil
	}
	return d.(*Daenerys)
}

type appKeyType struct{}

var appkey = appKeyType{}

func WithAPPKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, appkey, key)
}

func FromContext(ctx context.Context) (*Daenerys, bool) {
	key := ctx.Value(appkey)
	if key == nil {
		return nil, false
	}
	if d := GlobalNamespace.Get(key.(string)); d != nil {
		return d, true
	}
	return nil, false
}

func init() {
	log.RegisteCtx(func(ctx context.Context) (string, string) {
		return "namespace", ns.GetNamespace(ctx)
	})
}

func InitNamespace(dir string) {
	fd, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	fds, err := fd.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, finfo := range fds {
		if !finfo.IsDir() {
			continue
		}
		GlobalNamespace.Add(
			finfo.Name(),
			Namespace(finfo.Name()),
			ConfigPath(path.Join(dir, finfo.Name(), "config.toml")),
			App(Default.App),
			Name(Default.Name),
			Version(Default.Version),
			Deps(Default.Deps),
			Kit(Default.Kit),
		)
	}
}
