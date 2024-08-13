package events

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/Linxhhh/LinInk/api/proto/feed"
	samarax "github.com/Linxhhh/LinInk/pkg/saramax"
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

// -------------------------------------------- Article Publish Event Consumer -----------------------------------------------------------------

type ArticlePublishEventConsumer struct {
	client sarama.Client
	svc    feed.FeedServiceClient
}

func NewArticlePublishEventConsumer(client sarama.Client, svc feed.FeedServiceClient) *ArticlePublishEventConsumer {
	return &ArticlePublishEventConsumer{
		svc:    svc,
		client: client,
	}
}

// Start 启动 goroutine 消费事件
func (r *ArticlePublishEventConsumer) Start(consumerGroup string) error {

	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, r.client)
	if err != nil {
		return err
	}
	log.Println("start consumer ", consumerGroup)

	go func() {
		err := cg.Consume(context.Background(), []string{TopicArticlePublishEvent}, samarax.NewConsumer[ArticlePublishEvent](r.consume))
		if err != nil {
			log.Println("退出了消费循环异常", err)
		}
	}()
	return err
}

// Consume 消费事件
func (r *ArticlePublishEventConsumer) consume(msg *sarama.ConsumerMessage, evt ArticlePublishEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ext := map[string]string{
		"uid":   strconv.FormatInt(evt.Uid, 10),
		"aid":   strconv.FormatInt(evt.Aid, 10),
		"title": evt.Title,
	}
	val, err := json.Marshal(ext)
	if err != nil {
		return err
	}

	_, err = r.svc.Create(ctx, &feed.CreateRequest{
		Feed: &feed.FeedEvent{
			Type: TopicArticlePublishEvent,
			Ext:  string(val),
		},
	})
	return err
}