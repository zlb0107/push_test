package client

import (
	"net"
)

type dialer interface {
	Dial(host string) (socket, error)
}

type defaultDialer struct {
	opts Options
}

func (d defaultDialer) Dial(host string) (socket, error) {
	conn, err := net.DialTimeout("tcp4", host, d.opts.DialTimeout)
	if err != nil {
		return nil, err
	}
	socket := newIKSocket(d.opts.Kit.G(), conn, d.opts.CallOptions.RequestTimeout)
	return socket, nil
}
