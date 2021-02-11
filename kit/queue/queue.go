package queue

import (
	"context"
)

type Queue interface {
	Open(ctx context.Context, clientId string) error
	Close() error
	Publish(ctx context.Context, topic string, msg *Message) error
	Subscribe(topic string, receiverChan chan<- []byte) error
	PublishAtMostOnce(ctx context.Context, topic string, msg *Message) error
	SubscribeAtMostOnce(topic string, receiverChan chan<- []byte) error
}

