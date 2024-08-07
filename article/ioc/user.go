package ioc

import (
	"github.com/Linxhhh/LinInk/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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