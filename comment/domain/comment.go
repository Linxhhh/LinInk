package domain

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id       int64         `json:"id"`
	Aid      int64         `json:"aid"`
	Uid      int64         `json:"uid"`
	Content  string        `json:"content"`
	RootId   sql.NullInt64 `json:"rootId"`   // 根评论
	ParentId sql.NullInt64 `json:"parentId"` // 父评论
	Ctime    time.Time     `json:"ctime"`

	// 应该从 article 服务中获取
	AuthorName string `json:"authorName"`
	ParentName string `json:"parentName"`
}
