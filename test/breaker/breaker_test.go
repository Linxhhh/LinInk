package breaker

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Linxhhh/LinInk/api/proto/user"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BreakerTestSuite struct {
	suite.Suite
	cli *clientv3.Client
}

func (s *BreakerTestSuite) SetupSuite() {
	cli, err := clientv3.NewFromURL("localhost:12379")
	require.NoError(s.T(), err)
	s.cli = cli
}

func (s *BreakerTestSuite) TestServer() {
	go func() {
		s.startServer(&Server{}, "localhost:3333")
	}()
	go func() {
		s.startServer(&Server{}, "localhost:3334")
	}()
	s.startServer(&FailServer{}, "localhost:3335")
}

// start a server
func (s *BreakerTestSuite) startServer(svr user.UserServiceServer, addr string) {

	// init breaker interceptor
	breakerInterceptorBuilder := NewInterceptorBuilder("test", 10 * time.Second, 10)
	breakerInterceptor := breakerInterceptorBuilder.BuildServerInterceptor()

	// init grpc server
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(breakerInterceptor))
	user.RegisterUserServiceServer(grpcServer, svr)

	// listen port
	listener, err := net.Listen("tcp", addr)
	require.NoError(s.T(), err)

	// init etcd client
	client, err := clientv3.NewFromURL("localhost:12379")
	require.NoError(s.T(), err)

	// 节点管理模块
	manager, err := endpoints.NewManager(client, "service/user")
	require.NoError(s.T(), err)

	// 创建一个 15s 租约
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lease, err := client.Grant(ctx, 15)
	require.NoError(s.T(), err)

	// 注册节点，并携带租约
	key := fmt.Sprintf("%s/%s", "service/user", addr)
	ep := endpoints.Endpoint{Addr: addr}
	err = manager.AddEndpoint(ctx, key, ep, clientv3.WithLease(lease.ID))
	require.NoError(s.T(), err)
	log.Println("注册节点：" + addr)

	// 开启续约
	ch, err := client.KeepAlive(context.Background(), lease.ID)
	require.NoError(s.T(), err)

	go func() {
		// 当 cancel 被调用时，会退出这个循环
		for chResp := range ch {
			log.Println("续约 resp: ", chResp.String())
		}
	}()

	// start serve
	err = grpcServer.Serve(listener)
	require.NoError(s.T(), err)
}


func (s *BreakerTestSuite) TestClient() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	require.NoError(s.T(), err)
	cc, err := grpc.Dial(
		"etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`
{
  "loadBalancingConfig": [{"round_robin": {}}],
  "methodConfig":  [
    {
      "name": [{"service":  "user.UserService"}],
      "retryPolicy": {
        "maxAttempts": 4,
        "initialBackoff": "0.01s",
        "maxBackoff": "0.1s",
        "backoffMultiplier": 2.0,
        "retryableStatusCodes": ["UNAVAILABLE"]
      }
    }
  ]
}
`))
	require.NoError(t, err)
	client := user.NewUserServiceClient(cc)
	for i := 0; i < 20; i++ {
		resp, err := client.Profile(context.TODO(), &user.ProfileRequest{})
		require.NoError(t, err)
		t.Log(resp.User)
	}
}


func TestBreaker(t *testing.T) {
	suite.Run(t, new(BreakerTestSuite))
}

/*
	测试结果显示，在开启 Failover 的情况下，仍然触发了熔断（默认 5 次调用错误触发）
	说明 Failover 机制与熔断机制互不冲突
*/