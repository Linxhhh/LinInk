package ioc

import (
	"github.com/IBM/sarama"
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