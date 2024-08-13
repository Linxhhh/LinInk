package ioc

import (
	"log"

	"github.com/Linxhhh/LinInk/api/proto/article"
	"github.com/Linxhhh/LinInk/api/proto/code"
	"github.com/Linxhhh/LinInk/api/proto/follow"
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

func InitCodeRpcClient(cli *clientv3.Client) code.CodeServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/codeService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := code.NewCodeServiceClient(conn)
	return client
}


func InitArticleRpcClient(cli *clientv3.Client) article.ArticleServiceClient {

	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, _ := grpc.NewClient(
		"etcd:///LinInk/articleService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	client := article.NewArticleServiceClient(conn)
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