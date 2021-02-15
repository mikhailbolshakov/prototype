package queue

import (
	"context"
)

type QueueType int

const (
		QUEUE_TYPE_AT_LEAST_ONCE = iota
		QUEUE_TYPE_AT_MOST_ONCE
	)

type Queue interface {
	Open(ctx context.Context, clientId string) error
	Close() error
	Publish(ctx context.Context, qt QueueType, topic string, msg *Message) error
	Subscribe(qt QueueType, topic string, receiverChan chan<- []byte) error
}

