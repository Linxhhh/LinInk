package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Linxhhh/LinInk/feed/domain"
	"github.com/Linxhhh/LinInk/feed/repository/cache"
	"github.com/Linxhhh/LinInk/feed/repository/dao"
)

var FolloweesNotFound = cache.FolloweesNotFound

type FeedRepository interface {
	// 推事件
	CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error
	FindPushEvents(ctx context.Context, uid int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error)
	FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error)

	// 拉事件
	CreatePullEvent(ctx context.Context, event domain.FeedEvent) error
	FindPullEvents(ctx context.Context, uids []int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error)
	FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error)
}

type feedEventRepo struct {
	pullDao   dao.FeedPullEventDAO
	pushDao   dao.FeedPushEventDAO
	feedCache cache.FeedEventCache
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO, feedCache cache.FeedEventCache) FeedRepository {
	return &feedEventRepo{
		pullDao:   pullDao,
		pushDao:   pushDao,
		feedCache: feedCache,
	}
}

// --------------------------------------------------------- 推事件 -------------------------------------------------------------------

func (f *feedEventRepo) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	pushEvents := make([]dao.FeedPushEvent, 0, len(events))
	for _, e := range events {
		pushEvents = append(pushEvents, convertToPushEventDao(e))
	}
	return f.pushDao.CreatePushEvents(ctx, pushEvents)
}

func (f *feedEventRepo) FindPushEvents(ctx context.Context, uid int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.GetPushEvents(ctx, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.GetPushEventsWithTyp(ctx, typ, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

// --------------------------------------------------------- 拉事件 -------------------------------------------------------------------

func (f *feedEventRepo) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, convertToPullEventDao(event))
}

func (f *feedEventRepo) FindPullEvents(ctx context.Context, uids []int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEvents(ctx, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp time.Time, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventListWithTyp(ctx, typ, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

// ------------------------------------------------------- 辅助函数 -------------------------------------------------------------------


func (f *feedEventRepo) SetFollowees(ctx context.Context, follower int64, followees []int64) error {
	return f.feedCache.SetFollowees(ctx, follower, followees)
}

func (f *feedEventRepo) GetFollowees(ctx context.Context, follower int64) ([]int64, error) {
	followees, err := f.feedCache.GetFollowees(ctx, follower)
	if errors.Is(err, cache.FolloweesNotFound) {
		return nil, FolloweesNotFound
	}
	return followees, err
}

func convertToPushEventDao(event domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
		Id:      event.Id,
		Uid:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}
}

func convertToPullEventDao(event domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Id:      event.Id,
		Uid:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}

}

func convertToPushEventDomain(event dao.FeedPushEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		Id:    event.Id,
		Uid:   event.Uid,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}

func convertToPullEventDomain(event dao.FeedPullEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		Id:    event.Id,
		Uid:   event.Uid,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}