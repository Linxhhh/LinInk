package ioc

import (
	"log"

	"github.com/Linxhhh/LinInk/api/proto/user"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserRpcClient(cli *clientv3.Client) user.UserServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/userService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := user.NewUserServiceClient(conn)
	return client
}