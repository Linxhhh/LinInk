package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/user"
	"github.com/Linxhhh/LinInk/user/domain"
	"github.com/Linxhhh/LinInk/user/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserServiceServer(svc *service.UserService) *UserServiceServer {
	return &UserServiceServer{
		svc: svc,
	}
}

func (server *UserServiceServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	err := server.svc.SignUp(ctx, convertToDomain(req.GetUser()))
	return &pb.SignUpResponse{}, err
}

func (server *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	return &pb.LoginResponse{User: convertToResp(user)}, err
}

func (server *UserServiceServer) Edit(ctx context.Context, req *pb.EditInfoRequest) (*pb.EditInfoResponse, error) {
	err := server.svc.Edit(ctx, convertToDomain(req.GetUser()))
	return &pb.EditInfoResponse{}, err
}

func (server *UserServiceServer) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	user, err := server.svc.Profile(ctx, req.GetUid())
	return &pb.ProfileResponse{User: convertToResp(user)}, err
}

func (server *UserServiceServer) FindOrCreate(ctx context.Context, req *pb.FindOrCreateRequest) (*pb.FindOrCreateResponse, error) {
	uid, err := server.svc.FindOrCreate(ctx, req.GetPhone())
	return &pb.FindOrCreateResponse{Uid: uid}, err
}


// 类型转换：pb.User -> domain.User
func convertToDomain(u *pb.User) domain.User {
	domainUser := domain.User{}
	if u != nil {
		domainUser.Id = u.GetId()
		domainUser.Email = u.GetEmail()
		domainUser.Phone = u.GetPhone()
		domainUser.Password = u.GetPassword()
		domainUser.NickName = u.GetNickname()
		domainUser.Introduction = u.GetIntroduction()
		domainUser.Birthday = u.GetBirthday().AsTime()
	}
	return domainUser
}

// 类型转换：domain.User -> pb.User
func convertToResp(u domain.User) *pb.User {
	resp := &pb.User{
		Id:           u.Id,
		Email:        u.Email,
		Phone:        u.Phone,
		Password:     u.Password,
		Nickname:     u.NickName,
		Introduction: u.Introduction,
		Birthday:     timestamppb.New(u.Birthday),
	}
	return resp
}
