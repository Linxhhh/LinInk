package grpc

import (
	"context"
	"encoding/json"

	pb "github.com/Linxhhh/LinInk/api/proto/feed"
	"github.com/Linxhhh/LinInk/feed/domain"
	"github.com/Linxhhh/LinInk/feed/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedServiceServer struct {
	pb.UnimplementedFeedServiceServer
	svc *service.FeedEventService
}

func NewFeedServiceServer(svc *service.FeedEventService) *FeedServiceServer {
	return &FeedServiceServer{
		svc: svc,
	}
}

func (server *FeedServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	err := server.svc.CreateFeedEvent(ctx, convertToDomian(req.GetFeed()))
	return &pb.CreateResponse{}, err
}

func (server *FeedServiceServer) GetList(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	pullEvts, pushEvts, err := server.svc.GetFeedEventList(ctx, req.GetUid(), req.GetPushEvtTimestamp().AsTime(), req.GetPullEvtTimestamp().AsTime(), req.GetLimit())
	if err != nil {
		return &pb.GetListResponse{}, err
	}
	var pullEvtlist, pushEvtlist []*pb.FeedEvent
	for _, evt := range pullEvts {
		pullEvtlist = append(pullEvtlist, convertToPb(evt))
	}
	for _, evt := range pushEvts {
		pushEvtlist = append(pushEvtlist, convertToPb(evt))
	}
	return &pb.GetListResponse{PullEvtList: pullEvtlist, PushEvtList: pushEvtlist}, err
}

// 类型转换：pb.FeedEvent -> domain.FeedEvent
func convertToDomian(f *pb.FeedEvent) domain.FeedEvent {
	ext := map[string]string{}
	_ = json.Unmarshal([]byte(f.Ext), &ext)
	return domain.FeedEvent{
		Id:    f.GetId(),
		Uid:   f.GetUid(),
		Type:  f.GetType(),
		Ctime: f.GetCtime().AsTime(),
		Ext:   ext,
	}
}

// 类型转换：domain.FeedEvent -> pb.FeedEvent
func convertToPb(f domain.FeedEvent) *pb.FeedEvent {
	ext, _ := json.Marshal(f.Ext)
	return &pb.FeedEvent{
		Id:    f.Id,
		Uid:   f.Uid,
		Type:  f.Type,
		Ctime: timestamppb.New(f.Ctime),
		Ext:   string(ext),
	}
}
