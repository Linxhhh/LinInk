package service

import (
	"context"
	"fmt"

	"github.com/Linxhhh/LinInk/api/proto/user"
	"github.com/Linxhhh/LinInk/follow/domain"
	"github.com/Linxhhh/LinInk/follow/repository"
)

type FollowService struct {
	repo    repository.FollowRepository
	userCli user.UserServiceClient
}

func NewFollowService(repo repository.FollowRepository, userCli user.UserServiceClient) *FollowService {
	return &FollowService{
		repo:    repo,
		userCli: userCli,
	}
}

func (svc *FollowService) Follow(ctx context.Context, follower_id, followee_id int64) error {
	return svc.repo.Follow(ctx, follower_id, followee_id)
}

func (svc *FollowService) CancelFollow(ctx context.Context, follower_id, followee_id int64) error {
	return svc.repo.CancelFollow(ctx, follower_id, followee_id)
}

func (svc *FollowService) GetFollowData(ctx context.Context, uid int64) (domain.FollowData, error) {
	return svc.repo.GetFollowData(ctx, uid)
}

func (svc *FollowService) GetFollowed(ctx context.Context, follower_id, followee_id int64) (bool, error) {
	return svc.repo.GetFollowed(ctx, follower_id, followee_id)
}

func (svc *FollowService) GetFolloweeList(ctx context.Context, follower_id int64, page, pageSize int) ([]domain.FollowListData, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	followRelations, err := svc.repo.GetFolloweeList(ctx, follower_id, limit, offset)
	if err != nil {
		return nil, err
	}

	// 返回 followee List
	var followeeList []domain.FollowListData
	for _, f := range followRelations {

		// 获取用户信息
		userInfo, _ := svc.userCli.Profile(ctx, &user.ProfileRequest{Uid: f.Followee})

		// 获取关注信息
		followData, _ := svc.repo.GetFollowData(ctx, f.Followee)

		var followListItem = domain.FollowListData{
			Uid:        f.Followee,
			NickName:   userInfo.User.Nickname,
			Followees:  followData.Followees,
			Followers:  followData.Followers,
			IsFollowed: true,
		}
		followeeList = append(followeeList, followListItem)
	}
	return followeeList, err
}

func (svc *FollowService) GetFollowerList(ctx context.Context, followee_id int64, page, pageSize int) ([]domain.FollowListData, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	followRelations, err := svc.repo.GetFollowerList(ctx, followee_id, limit, offset)
	if err != nil {
		return nil, err
	}

	// 返回 follower List
	var followerList []domain.FollowListData
	for _, f := range followRelations {

		// 获取用户信息
		userInfo, _ := svc.userCli.Profile(ctx, &user.ProfileRequest{Uid: f.Follower})

		// 获取关注信息
		followData, _ := svc.repo.GetFollowData(ctx, f.Follower)

		// 是否已经关注
		isFollowed, err := svc.repo.GetFollowed(ctx, followee_id, f.Follower)
		if err != nil {
			return []domain.FollowListData{}, fmt.Errorf("“是否关注”数据获取失败！")
		}

		var followListItem = domain.FollowListData{
			Uid:        f.Follower,
			NickName:   userInfo.User.Nickname,
			Followees:  followData.Followees,
			Followers:  followData.Followers,
			IsFollowed: isFollowed,
		}
		followerList = append(followerList, followListItem)
	}
	return followerList, err
}

func (svc *FollowService) GetFolloweeIdList(ctx context.Context, follower_id int64, page, pageSize int) ([]int64, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	followRelations, err := svc.repo.GetFolloweeList(ctx, follower_id, limit, offset)
	if err != nil {
		return nil, err
	}

	// 返回 followee List
	var followeeList []int64
	for _, f := range followRelations {
		followeeList = append(followeeList, f.Followee)
	}
	return followeeList, err
}

func (svc *FollowService) GetFollowerIdList(ctx context.Context, followee_id int64, page, pageSize int) ([]int64, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	followRelations, err := svc.repo.GetFollowerList(ctx, followee_id, limit, offset)
	if err != nil {
		return nil, err
	}

	// 返回 follower List
	var followerList []int64
	for _, f := range followRelations {
		followerList = append(followerList, f.Follower)
	}
	return followerList, err
}
