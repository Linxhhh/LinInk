package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

const (
	TopicUserSyncEvent = "user_sync_event"
)

type UserSyncEvent struct {
	Id           int64  `json:"id"`
	NickName     string `json:"nickName"`
	Introduction string `json:"introduction"`
}

// -------------------------------------------- Article Publish Sync Producer -----------------------------------------------------------------

type UserSyncEventProducer struct {
	producer sarama.SyncProducer
}

func NewUserSyncEventProducer(producer sarama.SyncProducer) *UserSyncEventProducer {
	return &UserSyncEventProducer{producer: producer}
}

// Produce 生产事件
func (s *UserSyncEventProducer) Produce(evt UserSyncEvent) error {

	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicUserSyncEvent,
		Value: sarama.StringEncoder(val),
	})
	return err
}