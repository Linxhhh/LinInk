package grpc

import (
	"context"
	"database/sql"

	pb "github.com/Linxhhh/LinInk/api/proto/comment"
	"github.com/Linxhhh/LinInk/comment/domain"
	"github.com/Linxhhh/LinInk/comment/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommentServiceServer struct {
	pb.UnimplementedCommentServiceServer
	svc *service.CommentService
}

func NewCommentServer(svc *service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{
		svc: svc,
	}
}

func (server *CommentServiceServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	err := server.svc.Create(ctx, changeToDomain(req.GetCmt()))
	return &pb.CreateResponse{}, err
}

func (server *CommentServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := server.svc.Delete(ctx, req.GetId())
	return &pb.DeleteResponse{}, err
}

func (server *CommentServiceServer) RootComment(ctx context.Context, req *pb.RootCommentRequest) (*pb.RootCommentResponse, error) {
	cmts, err := server.svc.RootComment(ctx, req.GetAid(), req.GetMinId(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &pb.RootCommentResponse{Cmts: changeToPb(cmts)}, nil
}

func (server *CommentServiceServer) ChildComment(ctx context.Context, req *pb.ChildCommentRequest) (*pb.ChildCommentResponse, error) {
	cmts, err := server.svc.ChildComment(ctx, req.GetAid(), req.GetRootId(), req.GetMaxId(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &pb.ChildCommentResponse{Cmts: changeToPb(cmts)}, nil
}

// changeToDomain 类型转换
func changeToDomain(cmt *pb.Comment) domain.Comment {
	return domain.Comment{
		Aid:     cmt.Aid,
		Uid:     cmt.Uid,
		Content: cmt.Content,
		RootId: sql.NullInt64{
			Int64: cmt.RootId,
			Valid: cmt.RootId > 0,
		},
		ParentId: sql.NullInt64{
			Int64: cmt.ParentId,
			Valid: cmt.ParentId > 0,
		},
	}
}

// changeToPb 类型转换
func changeToPb(cmts []domain.Comment) []*pb.Comment {
	pbCmts := make([]*pb.Comment, len(cmts))
	for i, cmt := range cmts {
		pbCmts[i] = &pb.Comment{
			Id:         cmt.Id,
			Aid:        cmt.Aid,
			Uid:        cmt.Uid,
			AuthorName: cmt.AuthorName,
			Content:    cmt.Content,
			RootId:     cmt.RootId.Int64,
			ParentId:   cmt.RootId.Int64,
			ParentName: cmt.ParentName,
			Ctime:      timestamppb.New(cmt.Ctime),
		}
	}
	return pbCmts
}
