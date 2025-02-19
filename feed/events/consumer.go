package events

type Consumer interface {
	Start(consumerGroup string) error
}