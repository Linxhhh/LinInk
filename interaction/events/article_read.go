package events

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/Linxhhh/LinInk/interaction/service"
	samarax "github.com/Linxhhh/LinInk/pkg/saramax"
)

const (
	TopicArticleReadEvent = "article_read_event"
)

type ArticleReadEvent struct {
	Aid int64
}

// -------------------------------------------- Article Read Event Consumer -----------------------------------------------------------------

type ArticleReadEventConsumer struct {
	client sarama.Client
	svc    *service.InteractionService
	biz    string
}

func NewArticleReadEventConsumer(client sarama.Client, svc *service.InteractionService) *ArticleReadEventConsumer {
	return &ArticleReadEventConsumer{
		svc:    svc,
		client: client,
		biz: "article",
	}
}

func (i *ArticleReadEventConsumer) Start(consumerGroup string) error {

	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, i.client)
	if err != nil {
		return err
	}

	go func() {
		er := cg.Consume(context.Background(), []string{TopicArticleReadEvent}, samarax.NewConsumer[ArticleReadEvent](i.consume))
		if er != nil {
			log.Print("退出消费", er)
		}
	}()
	return err
}

func (i *ArticleReadEventConsumer) consume(msg *sarama.ConsumerMessage, event ArticleReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.svc.IncrReadCnt(ctx, i.biz, event.Aid)
}

// -------------------------------------------------------- Batch Consume -----------------------------------------------------------------

func (i *ArticleReadEventConsumer) StartBatch(consumerGroup string) error {
	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, i.client)
	if err != nil {
		return err
	}
	log.Println("start consumer ", consumerGroup)

	go func() {
		er := cg.Consume(context.Background(), []string{TopicArticleReadEvent}, samarax.NewBatchConsumer[ArticleReadEvent](i.batchConsume))
		if er != nil {
			log.Print("退出消费", er)
		}
	}()
	return err
}

func (i *ArticleReadEventConsumer) batchConsume(msgs []*sarama.ConsumerMessage, events []ArticleReadEvent) error {

	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, evt.Aid)
	}

	// 链路超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return i.svc.BatchIncrReadCnt(ctx, bizs, bizIds)
}
