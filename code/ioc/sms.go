package ioc

import (
	"log"

	"github.com/Linxhhh/LinInk/api/proto/sms"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSmsRpcClient(cli *clientv3.Client) sms.SmsServiceClient {
	
	// Create resolver builder
	resolverBuilder, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("create resolver builder failed, err : %s", err)
	}

	// Conn gRPC
	conn, err := grpc.NewClient(
		"etcd:///LinInk/smsService",
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}

	client := sms.NewSmsServiceClient(conn)
	return client
}

/*
func InitSmsRpcClient() sms.SmsServiceClient {
	
	// Connect sms service server
	conn, err := grpc.Dial(
		"localhost:3335",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	client := sms.NewSmsServiceClient(conn)
	return client
}
*/