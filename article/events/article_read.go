package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

const (
	TopicArticleReadEvent = "article_read_event"
)

type ArticleReadEvent struct {
	Aid int64
}

// -------------------------------------------- Article Read Event Producer -----------------------------------------------------------------

type ArticleReadEventProducer struct {
	producer sarama.SyncProducer
}

func NewArticleReadEventProducer(producer sarama.SyncProducer) *ArticleReadEventProducer {
	return &ArticleReadEventProducer{producer: producer}
}

func (s *ArticleReadEventProducer) ProduceEvent(evt ArticleReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicArticleReadEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}