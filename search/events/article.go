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
	TopicArticleSyncEvent     = "article_sync_event"
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

type ArticleSyncEventConsumer struct {
	client sarama.Client
	svc    *service.SyncService
}

func NewArticleSyncEventConsumer(client sarama.Client, svc *service.SyncService) *ArticleSyncEventConsumer {
	return &ArticleSyncEventConsumer{
		svc:    svc,
		client: client,
	}
}

// Start 启动 goroutine 消费事件
func (r *ArticleSyncEventConsumer) Start(consumerGroup string) error {

	cg, err := sarama.NewConsumerGroupFromClient(consumerGroup, r.client)
	if err != nil {
		return err
	}
	log.Println("start consumer ", consumerGroup)

	go func() {
		err := cg.Consume(context.Background(), []string{TopicArticleSyncEvent}, samarax.NewConsumer[ArticleSyncEvent](r.consume))
		if err != nil {
			log.Println("退出了消费循环异常", err) // 首次运行，会退出，因为 topic 尚未创建
		}
	}()

	go func() {
		err := cg.Consume(context.Background(), []string{TopicArticleWithdrawEvent}, samarax.NewConsumer[ArticleWithdrawEvent](r.consumeWithdraw))
		if err != nil {
			log.Println("退出了消费循环异常", err) // 首次运行，会退出，因为 topic 尚未创建
		}
	}()
	return err
}

// Consume 消费事件
func (r *ArticleSyncEventConsumer) consume(msg *sarama.ConsumerMessage, evt ArticleSyncEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.svc.PutArticle(ctx, domain.Article{
		Id:       evt.Id,
		Title:    evt.Title,
		Content:  evt.Content,
		AuthorId: evt.AuthorId,
		Status:   evt.Status,
		Ctime:    evt.Ctime,
		Utime:    evt.Ctime,
	})
}

func (r *ArticleSyncEventConsumer) consumeWithdraw(msg *sarama.ConsumerMessage, evt ArticleWithdrawEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return r.svc.WithdrawArticle(ctx, evt.Id)
}
