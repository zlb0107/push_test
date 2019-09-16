package connect_pool

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/collinmsn/thrift-client-pool"
	"go_common_lib/connect_pool/gen-go/following"
	"net"
)

var Follow_con_pool *thrift_client_pool.ChannelClientPool

func init() {
	// create pool
	//10.111.13.16：9090
	hosts := []string{net.JoinHostPort("10.111.13.16", "9090")}
	//超时单位:ns
	Follow_con_pool = thrift_client_pool.NewChannelClientPool(100, 0, hosts, 20000000, 20000000,
		func(openedSocket thrift.TTransport) thrift_client_pool.Client {
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
			transportFactory := thrift.NewTBufferedTransportFactory(8192)
			//transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
			transport := transportFactory.GetTransport(openedSocket)
			return following.NewRelationServiceClientFactory(transport, protocolFactory)
		},
	)

}
