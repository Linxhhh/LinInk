package ioc

import (
	"github.com/Linxhhh/LinInk/api/proto/sms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
