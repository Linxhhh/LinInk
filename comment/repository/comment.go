package repository

import (
	"context"
	"time"

	"github.com/Linxhhh/LinInk/comment/domain"
	"github.com/Linxhhh/LinInk/comment/repository/dao"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, cmt domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	GetRootComment(ctx context.Context, aid, minID int64, limit int) ([]domain.Comment, error)
	GetChildComment(ctx context.Context, aid, rootId, maxId int64, limit int) ([]domain.Comment, error)

	GetParentUidById(ctx context.Context, id int64) (int64, error)
}

type commentRepository struct {
	dao dao.CommentDAO
}

func NewCommentRepository(dao dao.CommentDAO) CommentRepository {
	return &commentRepository{
		dao: dao,
	}
}

func (repo *commentRepository) CreateComment(ctx context.Context, cmt domain.Comment) error {
	return repo.dao.Insert(ctx, dao.Comment{
		Aid:      cmt.Aid,
		Uid:      cmt.Uid,
		Content:  cmt.Content,
		RootId:   cmt.RootId,
		ParentId: cmt.ParentId,
	})
}

func (repo *commentRepository) DeleteComment(ctx context.Context, id int64) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *commentRepository) GetRootComment(ctx context.Context, aid, minID int64, limit int) ([]domain.Comment, error) {
	cmts, err := repo.dao.FindParentComment(ctx, aid, minID, limit)
	if err != nil {
		return nil, err
	}
	return changeToDomain(cmts), nil
}

func (repo *commentRepository) GetChildComment(ctx context.Context, aid, rootId, maxId int64, limit int) ([]domain.Comment, error) {
	cmts, err := repo.dao.FindChildComment(ctx, aid, rootId, maxId, limit)
	if err != nil {
		return nil, err
	}
	return changeToDomain(cmts), nil
}

func (repo *commentRepository) GetParentUidById(ctx context.Context, id int64) (int64, error) {
	return repo.dao.GetUidById(ctx, id)
}

// cchangeToDomain 类型转换
func changeToDomain(cmts []dao.Comment) []domain.Comment {
	domainCmts := make([]domain.Comment, len(cmts))
	for i, cmt := range cmts {
		domainCmts[i] = domain.Comment{
			Id:       cmt.Id,
			Aid:      cmt.Aid,
			Uid:      cmt.Uid,
			Content:  cmt.Content,
			RootId:   cmt.RootId,
			ParentId: cmt.ParentId,
			Ctime:    time.UnixMilli(cmt.Ctime),
		}
	}
	return domainCmts
}
