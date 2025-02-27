package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Linxhhh/LinInk/api/proto/feed"
	"github.com/Linxhhh/LinInk/api/proto/interaction"
	"github.com/Linxhhh/LinInk/api/proto/user"
	"github.com/Linxhhh/LinInk/article/domain"
	"github.com/Linxhhh/LinInk/article/events"
	"github.com/Linxhhh/LinInk/article/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var ErrIncorrectArticleorAuthor = repository.ErrIncorrectArticleorAuthor

type ArticleService struct {
	repo       repository.ArticleRepository
	userCli    user.UserServiceClient
	interCli   interaction.InteractionServiceClient
	feedCli    feed.FeedServiceClient
	publishPdr *events.ArticlePublishEventProducer
	readPdr    *events.ArticleReadEventProducer
	syncPdr    *events.ArticleSyncEventProducer
	Biz        string
}

func NewArticleService(repo repository.ArticleRepository,
	userCli user.UserServiceClient, interCli interaction.InteractionServiceClient, feedCli feed.FeedServiceClient,
	publishPdr *events.ArticlePublishEventProducer, readPdr *events.ArticleReadEventProducer,
	syncPdr *events.ArticleSyncEventProducer) *ArticleService {

	return &ArticleService{
		repo:       repo,
		userCli:    userCli,
		interCli:   interCli,
		feedCli:    feedCli,
		publishPdr: publishPdr,
		readPdr:    readPdr,
		syncPdr:    syncPdr,
		Biz:        "article",
	}
}

func (as *ArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// create sub span
	var tracer = otel.Tracer("LinInk-article-service")
	_, subSpan := tracer.Start(ctx, "sub-span-articleservice-save", trace.WithAttributes(attribute.String("key1", "value1")))
	defer subSpan.End()

	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return art.Id, as.repo.Update(ctx, art)
	}
	return as.repo.Insert(ctx, art)
}

func (as *ArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	art, err := as.repo.Sync(ctx, art)
	if err != nil {
		return 0, err
	}

	var eg errgroup.Group

	// 生成 feed 事件
	eg.Go(func() error {
		return as.publishPdr.Produce(events.ArticlePublishEvent{
			Aid:   art.Id,
			Uid:   art.AuthorId,
			Title: art.Title,
		})
	})

	// 同步 article 到 elasticSearch
	eg.Go(func() error {
		return as.syncPdr.Produce(events.ArticleSyncEvent{
			Id:       art.Id,
			Title:    art.Title,
			Content:  art.Content,
			AuthorId: art.AuthorId,
			Status:   int32(art.Status),
			Utime:    art.Utime,
			Ctime:    art.Ctime,
		})
	})
	err = eg.Wait()
	return art.Id, err
}

func (as *ArticleService) Withdraw(ctx context.Context, uid int64, aid int64) error {
	err := as.repo.SyncStatus(ctx, uid, aid, domain.ArticleStatusPrivate)
	if err != nil {
		return err
	}

	// 同步 article status 到 elasticSearch
	err = as.syncPdr.ProduceWithdraw(events.ArticleWithdrawEvent{
		Id: aid,
	})
	return err
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

	// 添加 readCnt
	go func() {
		if err == nil {
			as.readPdr.ProduceEvent(events.ArticleReadEvent{Aid: aid})
		}
	}()

	// 获取 AuthorName
	user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: art.AuthorId})
	if err != nil {
		return domain.Article{}, errors.New("查找用户失败")
	}
	art.AuthorName = user.GetUser().GetNickname()
	return art, nil
}

func (as *ArticleService) PubList(ctx context.Context, startTime time.Time, limit int) ([]domain.Article, error) {
	arts, err := as.repo.GetPubList(ctx, startTime, limit)
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

func (as *ArticleService) CollectionList(ctx context.Context, uid int64, limit, offset int32) ([]domain.Article, error) {

	// 获取 aidList
	resp, err := as.interCli.CollectionList(ctx, &interaction.CollectionListRequest{
		Biz:    as.Biz,
		Uid:    uid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []domain.Article{}, err
	}

	// 获取 article List
	arts, err := as.repo.GetPubListByIdList(ctx, resp.GetAidList())
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

func (as *ArticleService) PubWorks(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error) {
	arts, err := as.repo.GetPubWorks(ctx, uid, limit, offset)
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

func (as *ArticleService) FeedList(ctx context.Context, uid int64, pushEvtTimestamp, pullEvtTimestamp time.Time, limit int64) ([]domain.Article, []domain.Article, error) {

	// 获取 Feed Event
	resp, err := as.feedCli.GetList(ctx, &feed.GetListRequest{
		Uid:              uid,
		PushEvtTimestamp: timestamppb.New(pushEvtTimestamp),
		PullEvtTimestamp: timestamppb.New(pullEvtTimestamp),
		Limit:            limit,
	})
	if err != nil {
		return []domain.Article{}, []domain.Article{}, err
	}

	pullEvtList := resp.GetPullEvtList()
	artPullList := make([]domain.Article, 0, len(pullEvtList))
	aidPullMap := make(map[int]bool, len(pullEvtList))
	for _, evt := range pullEvtList {
		
		aid := getAidFromEvt(evt)
		if aidPullMap[int(aid)] {
			continue
		}
		aidPullMap[int(aid)] = true

		art, err := as.repo.GetPubById(ctx, aid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return []domain.Article{}, []domain.Article{}, err
		}

		// 获取 AuthorName
		user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: art.AuthorId})
		if err != nil {
			return []domain.Article{}, []domain.Article{}, errors.New("查找用户失败")
		}
		art.AuthorName = user.GetUser().GetNickname()
		artPullList = append(artPullList, art)
	}

	pushEvtList := resp.GetPushEvtList()
	artPushList := make([]domain.Article, 0, len(pushEvtList))
	aidPushMap := make(map[int]bool, len(pushEvtList))
	for _, evt := range pushEvtList {

		aid := getAidFromEvt(evt)
		if aidPushMap[int(aid)] {
			continue
		}
		aidPushMap[int(aid)] = true

		art, err := as.repo.GetPubById(ctx, aid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return []domain.Article{}, []domain.Article{}, err
		}

		// 获取 AuthorName
		user, err := as.userCli.Profile(ctx, &user.ProfileRequest{Uid: art.AuthorId})
		if err != nil {
			return []domain.Article{}, []domain.Article{}, errors.New("查找用户失败")
		}
		art.AuthorName = user.GetUser().GetNickname()
		artPushList = append(artPushList, art)
	}

	return artPullList, artPushList, nil
}

func getAidFromEvt(evt *feed.FeedEvent) int64 {
	ext := map[string]string{}
	_ = json.Unmarshal([]byte(evt.GetExt()), &ext)
	aidStr, _ := ext["aid"]
	aid, _ := strconv.ParseInt(aidStr, 10, 64)
	return aid
}
