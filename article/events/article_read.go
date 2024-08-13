package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/Linxhhh/LinInk/api/proto/interaction"
	samarax "github.com/Linxhhh/LinInk/pkg/saramax"
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

// -------------------------------------------- Article Read Event Consumer -----------------------------------------------------------------

type ArticleReadEventConsumer struct {
	client sarama.Client
	svc    interaction.InteractionServiceClient
}

func NewArticleReadEventConsumer(client sarama.Client, svc interaction.InteractionServiceClient) *ArticleReadEventConsumer {
	return &ArticleReadEventConsumer{
		svc:    svc,
		client: client,
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
	_, err := i.svc.IncrReadCnt(ctx, &interaction.IncrReadCntRequest{
		Biz:   "article",
		BizId: event.Aid,
	})
	return err
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
	_, err := i.svc.BatchIncrReadCnt(ctx, &interaction.BatchIncrReadCntRequest{
		Bizs:   bizs,
		BizIds: bizIds,
	})
	return err
}
