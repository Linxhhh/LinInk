package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/code"
	"github.com/Linxhhh/LinInk/code/service"
)

type CodeServiceServer struct {
	pb.UnimplementedCodeServiceServer
	svc *service.CodeService
}

func NewCodeServiceServer(svc *service.CodeService) *CodeServiceServer {
	return &CodeServiceServer{
		svc: svc,
	}
}

func (server *CodeServiceServer) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	err := server.svc.Send(ctx, req.GetBiz(), req.GetPhone())
	return &pb.SendResponse{}, err
}

func (server *CodeServiceServer) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	err := server.svc.Verify(ctx, req.GetBiz(), req.GetPhone(), req.GetInputCode())
	return &pb.VerifyResponse{}, err
}