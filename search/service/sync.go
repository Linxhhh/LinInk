package service

import (
	"context"

	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/repository"
)

type SyncService struct {
	artSync  repository.ArticleRepository
	userSync repository.UserRepository
}

func NewSyncService(artSync repository.ArticleRepository, userSync repository.UserRepository) *SyncService {
	return &SyncService{
		artSync:  artSync,
		userSync: userSync,
	}
}

func (svc *SyncService) PutArticle(ctx context.Context, art domain.Article) error {
	return svc.artSync.PutArticle(ctx, art)
}

func (svc *SyncService) WithdrawArticle(ctx context.Context, id int64) error {
	return svc.artSync.WithdrawArticle(ctx, id)
}

func (svc *SyncService) PutUser(ctx context.Context, user domain.User) error {
	return svc.userSync.PutUser(ctx, user)
}