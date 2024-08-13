package failover

import (
	"context"

	"github.com/Linxhhh/LinInk/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// 模拟异常节点

type FailServer struct {
	user.UnimplementedUserServiceServer
}

var _ user.UserServiceServer = &FailServer{}

func (server *FailServer) Profile(ctx context.Context, req *user.ProfileRequest) (*user.ProfileResponse, error) {
	return &user.ProfileResponse{
		User: &user.User{
			Id:       1,
			Nickname: "异常节点",
		},
	}, status.Errorf(codes.Unavailable, "模拟服务端错误")
}

// 正常节点

type Server struct {
	user.UnimplementedUserServiceServer
}

var _ user.UserServiceServer = &Server{}

func (server *Server) Profile(ctx context.Context, req *user.ProfileRequest) (*user.ProfileResponse, error) {
	return &user.ProfileResponse{
		User: &user.User{
			Id:       2,
			Nickname: "正常节点",
		},
	}, nil
}