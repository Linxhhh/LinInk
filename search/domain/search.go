package domain

import "time"

type Article struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Status   int32     `json:"status"`
	AuthorId int64     `json:"authorId"`
	Ctime    time.Time `json:"ctime"`
	Utime    time.Time `json:"utime"`
}

type User struct {
	Id           int64  `json:"id"`
	NickName     string `json:"nickName"`
	Introduction string `json:"introduction"`
}

type SearchResult struct {
	Articles []Article
	Users    []User
}