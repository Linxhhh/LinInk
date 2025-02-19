package events

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/Linxhhh/LinInk/feed/domain"
	"github.com/Linxhhh/LinInk/feed/service"
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

type ArticlePublishEventConsumer struct {
	client sarama.Client
	svc    *service.FeedEventService
}

func NewArticlePublishEventConsumer(client sarama.Client, svc *service.FeedEventService) *ArticlePublishEventConsumer {
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
			log.Println("退出了消费循环异常", err)  // 首次运行，会退出，因为 topic 尚未创建
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

	return r.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: TopicArticlePublishEvent,
		Ext: ext,
	})
}