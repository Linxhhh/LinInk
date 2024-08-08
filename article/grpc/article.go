package grpc

import (
	"context"

	pb "github.com/Linxhhh/LinInk/api/proto/article"
	"github.com/Linxhhh/LinInk/article/domain"
	"github.com/Linxhhh/LinInk/article/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArticleServiceServer struct {
	pb.UnimplementedArticleServiceServer
	svc *service.ArticleService
}

func NewArticleServiceServer(svc *service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{
		svc: svc,
	}
}

func (server *ArticleServiceServer) Save(ctx context.Context, req *pb.SaveRequest) (*pb.SaveResponse, error) {
	aid, err := server.svc.Save(ctx, convertToDomain(req.GetArticle()))
	return &pb.SaveResponse{Aid: aid}, err
}

func (server *ArticleServiceServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	aid, err := server.svc.Publish(ctx, convertToDomain(req.GetArticle()))
	return &pb.PublishResponse{Aid: aid}, err
}

func (server *ArticleServiceServer) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	err := server.svc.Withdraw(ctx, req.GetUid(), req.GetAid())
	return &pb.WithdrawResponse{}, err
}

func (server *ArticleServiceServer) Count(ctx context.Context, req *pb.CountRequest) (*pb.CountResponse, error) {
	count, err := server.svc.Count(ctx, req.GetUid())
	return &pb.CountResponse{Count: count}, err
}

func (server *ArticleServiceServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	list, err := server.svc.List(ctx, req.GetUid(), int(req.GetPage()), int(req.GetPageSize()))
	listResp := []*pb.ArticleListElem{}
	for _, elem := range list {
		listResp = append(listResp, convertToPbList(elem))
	}
	return &pb.ListResponse{List: listResp}, err
}

func (server *ArticleServiceServer) Detail(ctx context.Context, req *pb.DetailRequest) (*pb.DetailResponse, error) {
	article, err := server.svc.Detail(ctx, req.GetUid(), req.GetAid())
	return &pb.DetailResponse{Article: convertToPb(article)}, err
}

func (server *ArticleServiceServer) PubDetail(ctx context.Context, req *pb.PubDetailRequest) (*pb.PubDetailResponse, error) {
	article, err := server.svc.PubDetail(ctx, req.GetAid())
	return &pb.PubDetailResponse{Article: convertToPb(article)}, err
}

func (server *ArticleServiceServer) PubList(ctx context.Context, req *pb.PubListRequest) (*pb.PubListResponse, error) {
	list, err := server.svc.PubList(ctx, req.GetTimestamp().AsTime(), int(req.GetLimit()), int(req.GetOffset()))
	listResp := []*pb.Article{}
	for _, elem := range list {
		listResp = append(listResp, convertToPb(elem))
	}
	return &pb.PubListResponse{List: listResp}, err
}

func (server *ArticleServiceServer) CollectionList(ctx context.Context, req *pb.CollectionListRequest) (*pb.CollectionListResponse, error) {
	list, err := server.svc.CollectionList(ctx, req.GetUid())
	listResp := []*pb.Article{}
	for _, elem := range list {
		listResp = append(listResp, convertToPb(elem))
	}
	return &pb.CollectionListResponse{List: listResp}, err
}

// 类型转换：pb.Article -> domain.Article
func convertToDomain(a *pb.Article) domain.Article {
	domainArticle := domain.Article{}
	if a != nil {
		domainArticle.Id = a.GetId()
		domainArticle.Title = a.GetTitle()
		domainArticle.Content = a.GetContent()
		domainArticle.AuthorId = a.GetAuthorId()
	}
	return domainArticle
}

// 类型转换：domain.Article -> pb.Article
func convertToPb(a domain.Article) *pb.Article {
	resp := &pb.Article{
		Id:         a.Id,
		Title:      a.Title,
		Content:    a.Content,
		AuthorId:   a.AuthorId,
		AuthorName: a.AuthorName,
		Status:     uint32(a.Status),
		Ctime:      timestamppb.New(a.Ctime),
		Utime:      timestamppb.New(a.Utime),
	}
	return resp
}

// 类型转换：pb.ArticleListElem -> domain.ArticleListElem
func convertToPbList(a domain.ArticleListElem) *pb.ArticleListElem {
	articleListElem := &pb.ArticleListElem{
		Id:       a.Id,
		Title:    a.Title,
		Abstract: a.Abstract,
		Status:   uint32(a.Status),
		Ctime:    timestamppb.New(a.Ctime),
		Utime:    timestamppb.New(a.Utime),
	}
	return articleListElem
}
