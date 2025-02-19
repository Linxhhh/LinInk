package events

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

const (
	TopicArticleSyncEvent = "article_sync_event"
	TopicArticleWithdrawEvent = "article_withdraw_event"
)

type ArticleSyncEvent struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Status   int32     `json:"status"`
	AuthorId int64     `json:"authorId"`
	Ctime    time.Time `json:"ctime"`
	Utime    time.Time `json:"utime"`
}

type ArticleWithdrawEvent struct {
	Id int64 `json:"id"`
}

// -------------------------------------------- Article Publish Sync Producer -----------------------------------------------------------------

type ArticleSyncEventProducer struct {
	producer sarama.SyncProducer
}

func NewArticleSyncEventProducer(producer sarama.SyncProducer) *ArticleSyncEventProducer {
	return &ArticleSyncEventProducer{producer: producer}
}

// Produce 生产事件
func (s *ArticleSyncEventProducer) Produce(evt ArticleSyncEvent) error {

	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicArticleSyncEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}

func (s *ArticleSyncEventProducer) ProduceWithdraw(evt ArticleWithdrawEvent) error {

	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicArticleWithdrawEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}