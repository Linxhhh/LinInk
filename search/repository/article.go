package repository

import (
	"context"

	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/repository/dao"
)

type ArticleRepository interface {
	PutArticle(ctx context.Context, art domain.Article) error
	WithdrawArticle(ctx context.Context, id int64) error
	SearchArticle(ctx context.Context, expression string) ([]domain.Article, error)
}

type articleRepository struct {
	artDAO dao.ArticleElasticDAO
}

func NewArticleRepository(artDAO dao.ArticleElasticDAO) ArticleRepository {
	return &articleRepository{
		artDAO: artDAO,
	}
}

func (repo *articleRepository) PutArticle(ctx context.Context, art domain.Article) error {
	return repo.artDAO.InputArticle(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		Status:   art.Status,
		AuthorId: art.AuthorId,
		Ctime:    art.Ctime,
		Utime:    art.Utime,
	})
}

func (repo *articleRepository) WithdrawArticle(ctx context.Context, id int64) error {
	return repo.artDAO.WithdrawArticle(ctx, id)
}

func (repo *articleRepository) SearchArticle(ctx context.Context, expression string) ([]domain.Article, error) {
	arts, err := repo.artDAO.SearchArticle(ctx, expression)
	if err != nil {
		return nil, err
	}
	domainArts := make([]domain.Article, 0, len(arts))
	for _, art := range arts {
		domainArt := domain.Article{
			Id:       art.Id,
			Title:    art.Title,
			Content:  art.Content,
			Status:   art.Status,
			AuthorId: art.AuthorId,
			Ctime:    art.Ctime,
			Utime:    art.Utime,
		}
		domainArts = append(domainArts, domainArt)
	}
	return domainArts, nil
}
