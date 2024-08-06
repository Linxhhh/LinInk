package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/sms"
	"github.com/Linxhhh/LinInk/sms/service"
)

type SmsServiceServer struct {
	pb.UnimplementedSmsServiceServer
	svc service.Sms
}

func NewSmsServiceServer(svc service.Sms) *SmsServiceServer {
	return &SmsServiceServer{
		svc: svc,
	}
}

func (server *SmsServiceServer) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	err := server.svc.Send(ctx, req.GetTplId(), req.GetArgs(), req.Numbers...)
	return &pb.SendResponse{}, err
}