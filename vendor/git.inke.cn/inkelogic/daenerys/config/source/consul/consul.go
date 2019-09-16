package consul

import (
	"fmt"
	"git.inke.cn/inkelogic/daenerys/config/source"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
)

// Currently a single consul reader
type consul struct {
	prefix      string
	stripPrefix string
	addr        string
	opts        source.Options
	client      *api.Client
}

var (
	// DefaultPrefix is the prefix that consul keys will be assumed to have if you
	// haven't specified one
	DefaultPrefix = "/service_config"
	DefaultHost   = "127.0.0.1:8500"
)

func (c *consul) Read() (*source.ChangeSet, error) {
	kv, _, err := c.client.KV().List(c.prefix, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil || len(kv) == 0 {
		return nil, fmt.Errorf("source not found: %s", c.prefix)
	}

	//to map[string]interface{}
	data, err := makeMap(c.opts.Encoder, kv, c.stripPrefix)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	//to byte
	b, err := c.opts.Encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    c.opts.Encoder.String(),
		Source:    c.String(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (c *consul) String() string {
	return "consul"
}

func (c *consul) Watch() (source.Watcher, error) {
	w, err := newWatcher(c.prefix, c.addr, c.String(), c.stripPrefix, c.opts.Encoder)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	// use default config
	config := api.DefaultConfig()

	a, ok := options.Context.Value(addressKey{}).(string)
	if ok {
		addr, port, err := net.SplitHostPort(a)
		if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		} else {
			config.Address = DefaultHost
		}
	}

	// create the client
	client, _ := api.NewClient(config)

	prefix := DefaultPrefix
	sp := ""
	f, ok := options.Context.Value(prefixKey{}).(string)
	if ok {
		prefix = f
	}

	if b, ok := options.Context.Value(stripPrefixKey{}).(bool); ok && b {
		sp = prefix
	}

	return &consul{
		prefix:      prefix,
		stripPrefix: sp,
		addr:        config.Address,
		opts:        options,
		client:      client,
	}
}
