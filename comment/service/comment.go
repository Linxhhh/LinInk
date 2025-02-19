package service

import (
	"context"
	"errors"
	"math"

	"github.com/Linxhhh/LinInk/api/proto/user"
	"github.com/Linxhhh/LinInk/comment/domain"
	"github.com/Linxhhh/LinInk/comment/repository"
)

type CommentService struct {
	repo    repository.CommentRepository
	userCli user.UserServiceClient
}

func NewCommentService(repo repository.CommentRepository, userCli user.UserServiceClient) *CommentService {
	return &CommentService{
		repo:    repo,
		userCli: userCli,
	}
}

func (svc *CommentService) Create(ctx context.Context, cmt domain.Comment) error {
	return svc.repo.CreateComment(ctx, cmt)
}

func (svc *CommentService) Delete(ctx context.Context, id int64) error {
	return svc.repo.DeleteComment(ctx, id)
}

func (svc *CommentService) RootComment(ctx context.Context, aid, minID int64, limit int) ([]domain.Comment, error) {

	if minID == 0 {
		minID = math.MaxInt64
	}

	cmts, err := svc.repo.GetRootComment(ctx, aid, minID, limit)
	if err != nil {
		return []domain.Comment{}, err
	}
	for i := range cmts {
		// 获取 AuthorName
		user, err := svc.userCli.Profile(ctx, &user.ProfileRequest{Uid: cmts[i].Uid})
		if err != nil {
			return []domain.Comment{}, errors.New("查找用户失败")
		}
		cmts[i].AuthorName = user.GetUser().GetNickname()
	}
	return cmts, nil
}

func (svc *CommentService) ChildComment(ctx context.Context, aid, rootId, maxId int64, limit int) ([]domain.Comment, error) {
	cmts, err := svc.repo.GetChildComment(ctx, aid, rootId, maxId, limit)
	if err != nil {
		return []domain.Comment{}, err
	}
	for i := range cmts {
		// 获取 AuthorName
		author, err := svc.userCli.Profile(ctx, &user.ProfileRequest{Uid: cmts[i].Uid})
		if err != nil {
			return []domain.Comment{}, errors.New("查找用户失败")
		}
		cmts[i].AuthorName = author.GetUser().GetNickname()

		// 获取 ParentName
		parentUid, _ := svc.repo.GetParentUidById(ctx, cmts[i].ParentId.Int64)
		parentUser, err := svc.userCli.Profile(ctx, &user.ProfileRequest{Uid: parentUid})
		if err != nil {
			return []domain.Comment{}, errors.New("查找用户失败")
		}
		cmts[i].ParentName = parentUser.GetUser().GetNickname()
	}
	return cmts, nil
}