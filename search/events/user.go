package events

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	samarax "github.com/Linxhhh/LinInk/pkg/saramax"
	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/service"
)

const (
	TopicUserSyncEvent = "user_sync_event"
)

type UserSyncEvent struct {
	Id           int64  `json:"id"`
	NickName     string `json:"nickName"`
	Introduction string `json:"introduction"`
}

type UserSyncEventConsumer struct {
	client sarama.Client
	svc    *service.SyncService
}

func NewUserSyncEventConsumer(client sarama.Client, svc *service.SyncService) *UserSyncEventConsumer {
	return &UserSyncEventConsumer{
		svc:    svc,
		client: client,
	}
}

// Start 启动 goroutine 消费事件
func (r *UserSyncEventConsumer) Start(consumerGroup string) error {

	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, r.client)
	if err != nil {
		return err
	}
	log.Println("start consumer ", consumerGroup)

	go func() {
		err := cg.Consume(context.Background(), []string{TopicUserSyncEvent}, samarax.NewConsumer[UserSyncEvent](r.consume))
		if err != nil {
			log.Println("退出了消费循环异常", err) // 首次运行，会退出，因为 topic 尚未创建
		}
	}()
	return err
}

// Consume 消费事件
func (r *UserSyncEventConsumer) consume(msg *sarama.ConsumerMessage, evt UserSyncEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.svc.PutUser(ctx, domain.User{
		Id:           evt.Id,
		NickName:     evt.NickName,
		Introduction: evt.Introduction,
	})
}
