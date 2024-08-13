package ioc

import (
	"github.com/IBM/sarama"
	"github.com/Linxhhh/LinInk/article/events"
)

func InitSaramaClient() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, cfg)
	if err != nil {
		panic(err)
	}
	
	return client
}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(apc *events.ArticlePublishEventConsumer, arc *events.ArticleReadEventConsumer) []events.Consumer {
	return []events.Consumer{apc, arc}
}
