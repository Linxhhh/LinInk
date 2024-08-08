package ioc

import (
	"github.com/Linxhhh/LinInk/api/proto/follow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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