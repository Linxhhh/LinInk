package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Linxhhh/LinInk/api/proto/follow"
	"github.com/Linxhhh/LinInk/feed/domain"
	"github.com/Linxhhh/LinInk/feed/repository"
	"golang.org/x/sync/errgroup"
)

type FeedEventService struct {
	repo    repository.FeedRepository
	follCli follow.FollowServiceClient
}

func NewFeedEventService(repo repository.FeedRepository, follCli follow.FollowServiceClient) *FeedEventService {
	return &FeedEventService{
		repo:    repo,
		follCli: follCli,
	}
}

func (f *FeedEventService) CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error {

	followee, err := feed.Ext.Get("uid")
	if err != nil {
		return err
	}
	uid, err := strconv.ParseInt(followee, 10, 64)
	if err != nil {
		return err
	}

	// 根据粉丝数量，判定是拉模型还是推模型
	resp, err := f.follCli.GetFollowData(ctx, &follow.GetFollowDataRequest{
		Uid: uid,
	})
	if err != nil {
		return err
	}

	if resp.GetFollowData().GetFollowers() > 100 {
		// 拉模型（等粉丝拉取）
		return f.repo.CreatePullEvent(ctx, domain.FeedEvent{
			Uid:  uid,
			Type: domain.ArticleFeedEvent,
			// Type:  feed.Type,
			Ctime: time.Now(),
			Ext:   feed.Ext,
		})
	} else {
		// 推模型（推送给粉丝）
		listResp, err := f.follCli.GetFollowerIdList(ctx, &follow.GetFollowerListRequest{
			FolloweeId: uid,
			Page:       1,
			PageSize:   100000,
		})
		if err != nil {
			return err
		}
		var events []domain.FeedEvent
		for _, follower := range listResp.GetFollowerList() {
			events = append(events, domain.FeedEvent{
				Uid:  follower,
				Type: domain.ArticleFeedEvent,
				// Type:  feed.Type,
				Ctime: time.Now(),
				Ext:   feed.Ext,
			})
		}
		return f.repo.CreatePushEvents(ctx, events)
	}
}

// GetFeedEventList 查询发件箱和收信箱
func (f *FeedEventService) GetFeedEventList(ctx context.Context, uid int64, pushEvtTimestamp, pullEvtTimestamp time.Time, limit int64) (pullEvents, pushEvents []domain.FeedEvent, err error) {

	var eg errgroup.Group
	pullEvents = make([]domain.FeedEvent, 0, limit)
	pushEvents = make([]domain.FeedEvent, 0, limit)


	eg.Go(func() error {
		// 获取关注列表
		listResp, err := f.follCli.GetFolloweeIdList(ctx, &follow.GetFolloweeListRequest{
			FollowerId: uid,
			Page:       1,
			PageSize:   100000,
		})
		if err != nil {
			return err
		}
		followees := listResp.GetFolloweeList()

		// 查询发件箱
		evts, err := f.repo.FindPullEvents(ctx, followees, pullEvtTimestamp, limit)
		if err != nil {
			return err
		}
		pullEvents = append(pullEvents, evts...)
		return nil
	})

	eg.Go(func() error {
		// 查询收件箱
		evts, err := f.repo.FindPushEvents(ctx, uid, pushEvtTimestamp, limit)
		if err != nil {
			return err
		}
		pushEvents = append(pushEvents, evts...)
		return nil
	})

	err = eg.Wait()
	return
}
