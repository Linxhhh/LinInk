package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

const (
	TopicArticlePublishEvent = "article_publish_event"
)

type ArticlePublishEvent struct {
	Uid   int64  // author_id
	Aid   int64  // article_id
	Title string // article_title
}

// -------------------------------------------- Article Publish Event Producer -----------------------------------------------------------------

type ArticlePublishEventProducer struct {
	producer sarama.SyncProducer
}

func NewArticlePublishEventProducer(producer sarama.SyncProducer) *ArticlePublishEventProducer {
	return &ArticlePublishEventProducer{producer: producer}
}

// Produce 生产事件
func (s *ArticlePublishEventProducer) Produce(evt ArticlePublishEvent) error {

	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicArticlePublishEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}