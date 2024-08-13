package samarax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type BatchConsumer[T any] struct {
	fn func(msg []*sarama.ConsumerMessage, event []T) error
}

func NewBatchConsumer[T any](fn func(msgs []*sarama.ConsumerMessage, t []T) error) *BatchConsumer[T] {
	return &BatchConsumer[T]{
		fn: fn,
	}
}

func (bc *BatchConsumer[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (bc *BatchConsumer[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (bc *BatchConsumer[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	
	msgsCh := claim.Messages()
	const batchSize = 10  // 批量操作，数量为 10

	for {
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 该批次超时，或者 consumer 被关闭，则不要求凑够一批
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					// channel 被关闭
					cancel()
					return nil
				}
				msgs = append(msgs, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					// 不中断，继续下一个
					session.MarkMessage(msg, "")
					continue
				}
				ts = append(ts, t)
			}
		}
		err := bc.fn(msgs, ts)
		if err == nil {
			for _, msg := range msgs {
				session.MarkMessage(msg, "")
			}
		} else {
			// 可以考虑重试，也可以在业务中重试
		}
		cancel()
	}
}