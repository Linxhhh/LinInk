package service

import (
	"context"
	"sync"

	"github.com/Linxhhh/LinInk/interaction/domain"
	"github.com/Linxhhh/LinInk/interaction/repository"
)

type InteractionService struct {
	repo    repository.InteractionRepository
}

func NewInteractionService(repo repository.InteractionRepository) *InteractionService {
	return &InteractionService{
		repo:    repo,
	}
}

func (svc *InteractionService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *InteractionService) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	return svc.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}

func (svc *InteractionService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.Like(ctx, biz, bizId, uid)
}

func (svc *InteractionService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.CancelLike(ctx, biz, bizId, uid)
}

func (svc *InteractionService) Collect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.Collect(ctx, biz, bizId, uid)
}

func (svc *InteractionService) CancelCollect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.CancelCollect(ctx, biz, bizId, uid)
}

/*
后续优化：分页查询
*/
func (svc *InteractionService) CollectionList(ctx context.Context, biz string, uid int64, limit, offset int) ([]int64, error) {
	return svc.repo.GetCollectionList(ctx, biz, uid, limit, offset)
}

func (svc *InteractionService) Share(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.Share(ctx, biz, bizId)
}

func (svc *InteractionService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interaction, error) {

	// 获取（阅读、点赞、收藏）数据
	i, err := svc.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interaction{}, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 是否已经点赞
	go func() {
		defer wg.Done()
		i.IsLiked, _ = svc.repo.GetLike(ctx, biz, bizId, uid)
	}()

	// 是否已经收藏
	go func() {
		defer wg.Done()
		i.IsCollected, _ = svc.repo.GetCollection(ctx, biz, bizId, uid)
	}()

	wg.Wait()

	return i, err
}

func (svc *InteractionService) BatchGet(ctx context.Context, biz string, bizIds []int64) (intrs []domain.Interaction, err error) {

	intrs = make([]domain.Interaction, len(bizIds))

	for i := 0; i < len(bizIds); i++ {
		intr, err := svc.repo.Get(ctx, biz, bizIds[i]); 
		if err != nil {
			return nil, err
		}
		intrs[i] = intr
	}

	return
}