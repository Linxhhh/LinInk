package ioc

import (
	"log"

	"github.com/Linxhhh/LinInk/api/proto/feed"
	"github.com/Linxhhh/LinInk/api/proto/interaction"
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

func InitFeedRpcClient(cli *clientv3.Client) feed.FeedServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/feedService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := feed.NewFeedServiceClient(conn)
	return client
}

func InitInteractionRpcClient(cli *clientv3.Client) interaction.InteractionServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/interactionService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := interaction.NewInteractionServiceClient(conn)
	return client
}

/*

// 使用 dns 域名解析来连接

func InitUserRpcClient() user.UserServiceClient {

	// Connect user service server
	conn, err := grpc.Dial(
		"localhost:3333",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	client := user.NewUserServiceClient(conn)
	return client
}

func InitInteractionRpcClient() interaction.InteractionServiceClient {

	// Connect interaction service server
	conn, err := grpc.Dial(
		"localhost:3337",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	client := interaction.NewInteractionServiceClient(conn)
	return client
}
*/
