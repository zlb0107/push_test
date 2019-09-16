package connect_pool

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/collinmsn/thrift-client-pool"
	"go_common_lib/connect_pool/gen-go/new_user_rank_service"
	"net"
)

var New_recall_con_pool *thrift_client_pool.ChannelClientPool

func init() {
	// create pool
	hosts := []string{net.JoinHostPort("10.111.95.228", "9091")}
	//超时单位:ns
	New_recall_con_pool = thrift_client_pool.NewChannelClientPool(100, 0, hosts, 20000000, 20000000,
		func(openedSocket thrift.TTransport) thrift_client_pool.Client {
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
			//transportFactory := thrift.NewTBufferedTransportFactory(8192)
			transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
			//transport,_ := transportFactory.GetTransport(openedSocket)
			transport := transportFactory.GetTransport(openedSocket)
			return new_user_rank_service.NewRecommendServiceClientFactory(transport, protocolFactory)
		},
	)

}
