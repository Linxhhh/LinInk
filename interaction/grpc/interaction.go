package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/interaction"
	"github.com/Linxhhh/LinInk/interaction/domain"
	"github.com/Linxhhh/LinInk/interaction/service"
)

type InteractionServiceServer struct {
	pb.UnimplementedInteractionServiceServer
	svc *service.InteractionService
}

func NewInteractionServiceServer(svc *service.InteractionService) *InteractionServiceServer {
	return &InteractionServiceServer{
		svc: svc,
	}
}

func (server *InteractionServiceServer) IncrReadCnt(ctx context.Context, req *pb.IncrReadCntRequest) (*pb.IncrReadCntResponse, error) {
	err := server.svc.IncrReadCnt(ctx, req.GetBiz(), req.GetBizId())
	return &pb.IncrReadCntResponse{}, err
}

func (server *InteractionServiceServer) BatchIncrReadCnt(ctx context.Context, req *pb.BatchIncrReadCntRequest) (*pb.BatchIncrReadCntResponse, error) {
	err := server.svc.BatchIncrReadCnt(ctx, req.GetBizs(), req.GetBizIds())
	return &pb.BatchIncrReadCntResponse{}, err
}

func (server *InteractionServiceServer) Like(ctx context.Context, req *pb.LikeRequest) (*pb.LikeResponse, error) {
	err := server.svc.Like(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &pb.LikeResponse{}, err
}

func (server *InteractionServiceServer) CancelLike(ctx context.Context, req *pb.CancelLikeRequest) (*pb.CancelLikeResponse, error) {
	err := server.svc.CancelLike(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &pb.CancelLikeResponse{}, err
}

func (server *InteractionServiceServer) Collect(ctx context.Context, req *pb.CollectRequest) (*pb.CollectResponse, error) {
	err := server.svc.Collect(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &pb.CollectResponse{}, err
}

func (server *InteractionServiceServer) CancelCollect(ctx context.Context, req *pb.CancelCollectRequest) (*pb.CancelCollectResponse, error) {
	err := server.svc.CancelCollect(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &pb.CancelCollectResponse{}, err
}

func (server *InteractionServiceServer) CollectionList(ctx context.Context, req *pb.CollectionListRequest) (*pb.CollectionListResponse, error) {
	list, err := server.svc.CollectionList(ctx, req.GetBiz(), req.GetUid(), int(req.GetLimit()), int(req.GetOffset()))
	return &pb.CollectionListResponse{AidList: list}, err
}

func (server *InteractionServiceServer) Share(ctx context.Context, req *pb.ShareRequest) (*pb.ShareResponse, error) {
	err := server.svc.Share(ctx, req.GetBiz(), req.GetBizId())
	return &pb.ShareResponse{}, err
}

func (server *InteractionServiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	interaction, err := server.svc.Get(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &pb.GetResponse{Interaction: convertToPb(interaction)}, err
}

func (server *InteractionServiceServer) BatchGet(ctx context.Context, req *pb.BatchGetRequest) (*pb.BatchGetResponse, error) {
	resp, err := server.svc.BatchGet(ctx, req.GetBiz(), req.GetBizIds())
	var interactions []*pb.Interaction
	for i := 0; i < len(resp); i++ {
		interactions = append(interactions, convertToPb(resp[i]))
	}
	return &pb.BatchGetResponse{Interactions: interactions}, err
}

// 类型转换：domain.Interaction -> pb.Interaction
func convertToPb(i domain.Interaction) *pb.Interaction {
	resp := &pb.Interaction{
		Id:          i.Id,
		Biz:         i.Biz,
		BizId:       i.BizId,
		ReadCnt:     i.ReadCnt,
		LikeCnt:     i.LikeCnt,
		CollectCnt:  i.CollectCnt,
		ShareCnt:    i.ShareCnt,
		IsLiked:     i.IsLiked,
		IsCollected: i.IsCollected,
	}
	return resp
}
