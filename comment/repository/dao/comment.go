package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type CommentDAO interface {
	Insert(ctx context.Context, cmt Comment) error
	Delete(ctx context.Context, id int64) error
	FindParentComment(ctx context.Context, aid, minID int64, limit int) ([]Comment, error)
	FindChildComment(ctx context.Context, aid, rootId, maxId int64, limit int) ([]Comment, error)

	GetUidById(ctx context.Context, id int64) (int64, error)
}

type GormCommentDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

// NewCommentDAO 新建一个数据库存储实例
func NewCommentDAO(m *gorm.DB, s []*gorm.DB) CommentDAO {
	return &GormCommentDAO{
		master: m,
		slaves: s,
	}
}

// RandSalve 随机获取从数据库
func (dao *GormCommentDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
	randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
	return randomSlave
}

// Insert 插入一条评论
func (dao *GormCommentDAO) Insert(ctx context.Context, cmt Comment) error {
	// 存储毫秒数
	cmt.Ctime = time.Now().UnixMilli()

	// 插入新记录
	return dao.master.WithContext(ctx).Create(&cmt).Error
}

// Delete 删除一条评论及其子评论
func (dao *GormCommentDAO) Delete(ctx context.Context, id int64) error {
	return dao.master.WithContext(ctx).Delete(&Comment{}, id).Error
}

// FindRootComment 获取根评论
func (dao *GormCommentDAO) FindParentComment(ctx context.Context, aid, minID int64, limit int) ([]Comment, error) {
	var res []Comment
	err := dao.RandSalve().WithContext(ctx).
		Where("aid = ? AND id < ? AND parent_id IS NULL", aid, minID).
		Order("id DESC").
		Limit(limit).
		Find(&res).Error
	return res, err
}

// FindChildComment 获取子评论
func (dao *GormCommentDAO) FindChildComment(ctx context.Context, aid, rootId, maxId int64, limit int) ([]Comment, error) {
	var res []Comment
	err := dao.RandSalve().WithContext(ctx).
		Where("aid = ? AND root_id = ? AND id > ?", aid, rootId, maxId).
		Order("id ASC").
		Limit(limit).
		Find(&res).Error
	return res, err
}

// Get 获取子评论
func (dao *GormCommentDAO) GetUidById(ctx context.Context, id int64) (int64, error) {
	var res Comment
	err := dao.RandSalve().WithContext(ctx).
		Where("id = ?", id).Find(&res).Error
	return res.Uid, err
}

type Comment struct {
	Id        int64 `gorm:"column:id;primaryKey"`
	Aid       int64 `gorm:"column:aid;index"`
	Uid       int64 `gorm:"column:uid;index"`
	Content   string
	RootId    sql.NullInt64 `gorm:"column:root_id;index"`
	ParentId  sql.NullInt64 `gorm:"column:parent_id;index"`
	ParentCmt *Comment      `gorm:"ForeignKey:ParentId;AssociationForeignKey:Id;constraint:OnDelete:CASCADE"`
	Ctime     int64

	// 为了方便，可以存储以下字段
	// AuthorName string `gorm:"column:author_name"`
	// ParentName string `gorm:"column:parent_name"`
}
