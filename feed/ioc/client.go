package ioc

import (
	"log"

	"github.com/Linxhhh/LinInk/api/proto/follow"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitFollowRpcClient(cli *clientv3.Client) follow.FollowServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/followService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := follow.NewFollowServiceClient(conn)
	return client
}

/*
func InitFollowRpcClient() follow.FollowServiceClient {

	// connect follow service server
	conn, err := grpc.Dial(
		"localhost:3338",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	client := follow.NewFollowServiceClient(conn)
	return client
}
*/