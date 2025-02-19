package service

import (
	"context"

	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/repository"
)

type SearchService struct {
	artRepo  repository.ArticleRepository
	userRepo repository.UserRepository
}

func NewSearchService(artRepo repository.ArticleRepository, userRepo repository.UserRepository) *SearchService {
	return &SearchService{
		artRepo:  artRepo,
		userRepo: userRepo,
	}
}

func (svc *SearchService) Search(ctx context.Context, expression string) (domain.SearchResult, error) {
	arts, err := svc.artRepo.SearchArticle(ctx, expression)
	if err != nil {
		return domain.SearchResult{}, err
	}
	users, err := svc.userRepo.SearchUser(ctx, expression)
	if err != nil {
		return domain.SearchResult{}, err
	}
	return domain.SearchResult{
		Articles: arts,
		Users:    users,
	}, nil
}
