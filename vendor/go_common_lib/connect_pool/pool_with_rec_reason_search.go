package connect_pool

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/collinmsn/thrift-client-pool"
	"net"
	"go_common_lib/connect_pool/gen-go/rec_rpc"
)

var RecReasonSearch_con_pool *thrift_client_pool.ChannelClientPool

func init() {
	// create pool
	hosts := []string{net.JoinHostPort("10.111.95.185", "18098")}
	//hosts := []string{net.JoinHostPort("127.0.0.1", "18098")}
	//超时单位:ns
	RecReasonSearch_con_pool = thrift_client_pool.NewChannelClientPool(1000, 0, hosts, 15000000, 15000000,
		func(openedSocket thrift.TTransport) thrift_client_pool.Client {
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
			transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
			//transport,_ := transportFactory.GetTransport(openedSocket)
			transport := transportFactory.GetTransport(openedSocket)
			return rec_rpc.NewRecTagSearchClientFactory(transport, protocolFactory)
		},
	)

}
