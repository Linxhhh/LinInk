package service

import (
	"context"
	"errors"
	"time"

	"github.com/Linxhhh/LinInk/api/proto/interaction"
	"github.com/Linxhhh/LinInk/api/proto/user"
	"github.com/Linxhhh/LinInk/article/domain"
	"github.com/Linxhhh/LinInk/article/repository"
)

var ErrIncorrectArticleorAuthor = repository.ErrIncorrectArticleorAuthor

type ArticleService struct {
	repo     repository.ArticleRepository
	userCli  user.UserServiceClient
	interCli interaction.InteractionServiceClient
	Biz      string
}

func NewArticleService(repo repository.ArticleRepository, userCli user.UserServiceClient, interCli interaction.InteractionServiceClient) *ArticleService {
	return &ArticleService{
		repo:     repo,
		userCli:  userCli,
		interCli: interCli,
		Biz:      "article",
	}
}

func (as *ArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return art.Id, as.repo.Update(ctx, art)
	}
	return as.repo.Insert(ctx, art)
}

func (as *ArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return as.repo.Sync(ctx, art)
}

func (as *ArticleService) Withdraw(ctx context.Context, uid int64, aid int64) error {
	return as.repo.SyncStatus(ctx, uid, aid, domain.ArticleStatusPrivate)
}

func (as *ArticleService) Count(ctx context.Context, uid int64) (int64, error) {
	return as.repo.CountByAuthor(ctx, uid)
}

func (as *ArticleService) List(ctx context.Context, uid int64, page, pageSize int) ([]domain.ArticleListElem, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return as.repo.GetListByAuthor(ctx, uid, offset, limit)
}

func (as *ArticleService) Detail(ctx context.Context, uid, aid int64) (domain.Article, error) {
	art, err := as.repo.GetById(ctx, aid)
	if err == nil && art.AuthorId != uid {
		return domain.Article{}, ErrIncorrectArticleorAuthor
	}
	return art, err
}

func (as *ArticleService) PubDetail(ctx context.Context, aid int64) (domain.Article, error) {
	art, err := as.repo.GetPubById(ctx, aid)
	if err != nil {
		return domain.Article{}, err
	}

	// 获取 AuthorName
	user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: art.AuthorId})
	if err != nil {
		return domain.Article{}, errors.New("查找用户失败")
	}
	art.AuthorName = user.GetUser().GetNickname()
	return art, nil
}

func (as *ArticleService) PubList(ctx context.Context, startTime time.Time, limit, offset int) ([]domain.Article, error) {
	arts, err := as.repo.GetPubList(ctx, startTime, offset, limit)
	if err != nil {
		return []domain.Article{}, err
	}
	for i := range arts {
		// 获取 AuthorName
		user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: arts[i].AuthorId})
		if err != nil {
			return []domain.Article{}, errors.New("查找用户失败")
		}
		arts[i].AuthorName = user.GetUser().GetNickname()
	}
	return arts, nil
}

func (as *ArticleService) CollectionList(ctx context.Context, uid int64) ([]domain.Article, error) {

	// 获取 aidList
	resp, err := as.interCli.CollectionList(ctx, &interaction.CollectionListRequest{
		Biz: as.Biz,
		Uid: uid,
	})
	if err != nil {
		return []domain.Article{}, err
	}

	// 获取 article List
	arts, err := as.repo.GetCollectionList(ctx, resp.GetAidList())
	if err != nil {
		return []domain.Article{}, err
	}
	for i := range arts {
		// 获取 AuthorName
		user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: arts[i].AuthorId})
		if err != nil {
			return []domain.Article{}, errors.New("查找用户失败")
		}
		arts[i].AuthorName = user.GetUser().GetNickname()
	}
	return arts, nil
}
