package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/follow"
	"github.com/Linxhhh/LinInk/follow/domain"
	"github.com/Linxhhh/LinInk/follow/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FollowServiceServer struct {
	pb.UnimplementedFollowServiceServer
	svc *service.FollowService
}

func NewFollowServiceServer(svc *service.FollowService) *FollowServiceServer {
	return &FollowServiceServer{
		svc: svc,
	}
}

func (server *FollowServiceServer) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error) {
	err := server.svc.Follow(ctx, req.GetFollowerId(), req.GetFolloweeId())
	return &pb.FollowResponse{}, err
}

func (server *FollowServiceServer) CancelFollow(ctx context.Context, req *pb.CancelFollowRequest) (*pb.CancelFollowResponse, error) {
	err := server.svc.CancelFollow(ctx, req.GetFollowerId(), req.GetFolloweeId())
	return &pb.CancelFollowResponse{}, err
}

func (server *FollowServiceServer) GetFollowData(ctx context.Context, req *pb.GetFollowDataRequest) (*pb.GetFollowDataResponse, error) {
	followdata, err := server.svc.GetFollowData(ctx, req.GetUid())
	return &pb.GetFollowDataResponse{FollowData: convertToPb(followdata)}, err
}

func (server *FollowServiceServer) GetFollowed(ctx context.Context, req *pb.GetFollowedRequest) (*pb.GetFollowedResponse, error) {
	isFollowed, err := server.svc.GetFollowed(ctx, req.GetFollowerId(), req.GetFolloweeId())
	return &pb.GetFollowedResponse{IsFollowed: isFollowed}, err
}

func (server *FollowServiceServer) GetFolloweeList(ctx context.Context, req *pb.GetFolloweeListRequest) (*pb.GetFolloweeListResponse, error) {
	followeeList, err := server.svc.GetFolloweeList(ctx, req.GetFollowerId(), int(req.GetPage()), int(req.GetPageSize()))
	return &pb.GetFolloweeListResponse{FolloweeList: followeeList}, err
}

func (server *FollowServiceServer) GetFollowerList(ctx context.Context, req *pb.GetFollowerListRequest) (*pb.GetFollowerListResponse, error) {
	followerList, err := server.svc.GetFollowerList(ctx, req.GetFolloweeId(), int(req.GetPage()), int(req.GetPageSize()))
	return &pb.GetFollowerListResponse{FollowerList: followerList}, err
}

// 类型转换：domain.FollowData -> pb.FollowData
func convertToPb(f domain.FollowData) *pb.FollowData {
	resp := &pb.FollowData{
		Id:    f.Id,
		Uid:   f.Uid,
		Followers: f.Followers,
		Followees: f.Followees,
		Ctime: timestamppb.New(f.Ctime),
		Utime: timestamppb.New(f.Utime),
	}
	return resp
}
