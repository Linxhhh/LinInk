package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/search"
	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SearchServiceServer struct {
	pb.SearchServiceServer
	syncSvc   *service.SyncService
	searchSvc *service.SearchService
}

func NewSearchServiceServer(syncSvc *service.SyncService, searchSvc *service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{
		syncSvc:   syncSvc,
		searchSvc: searchSvc,
	}
}

func (server *SearchServiceServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	res, err := server.searchSvc.Search(ctx, req.GetExpression())
	if err != nil {
		return &pb.SearchResponse{}, err
	}
	return &pb.SearchResponse{
		Articles: changeToPbArticles(res.Articles),
		Users:    changeToPbUser(res.Users),
	}, nil
}

func (server *SearchServiceServer) SyncArticle(ctx context.Context, req *pb.SyncArticleRequest) (*pb.SyncArticleResponse, error) {
	err := server.syncSvc.PutArticle(ctx, changeToDomainArticle(req.GetArticle()))
	return &pb.SyncArticleResponse{}, err
}

func (server *SearchServiceServer) SyncUser(ctx context.Context, req *pb.SyncUserRequest) (*pb.SyncUserResponse, error) {
	err := server.syncSvc.PutUser(ctx, changeToDomainUser(req.GetUser()))
	return &pb.SyncUserResponse{}, err
}

func changeToDomainArticle(art *pb.Article) domain.Article {
	return domain.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		Status:   art.Status,
		AuthorId: art.AuthorId,
		Ctime:    art.Ctime.AsTime(),
		Utime:    art.Utime.AsTime(),
	}
}

func changeToPbArticles(arts []domain.Article) []*pb.Article {
	pbArts := make([]*pb.Article, 0, len(arts))
	for _, art := range arts {
		pbArt := &pb.Article{
			Id:       art.Id,
			Title:    art.Title,
			Content:  art.Content,
			Status:   art.Status,
			AuthorId: art.AuthorId,
			Ctime:    timestamppb.New(art.Ctime),
			Utime:    timestamppb.New(art.Utime),
		}
		pbArts = append(pbArts, pbArt)
	}
	return pbArts
}

func changeToDomainUser(user *pb.User) domain.User {
	return domain.User{
		Id:           user.Id,
		NickName:     user.Nickname,
		Introduction: user.Introduction,
	}
}

func changeToPbUser(users []domain.User) []*pb.User {
	pbUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		pbUser := &pb.User{
			Id:           user.Id,
			Nickname:     user.NickName,
			Introduction: user.Introduction,
		}
		pbUsers = append(pbUsers, pbUser)
	}
	return pbUsers
}
