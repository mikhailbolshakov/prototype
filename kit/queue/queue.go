package queue

import (
	"context"
)

type QueueType int

func (q QueueType) String() string {
	if q == QUEUE_TYPE_AT_LEAST_ONCE {
		return "at-least-once"
	} else if  q == QUEUE_TYPE_AT_MOST_ONCE {
		return "at-most-once"
	}
	return ""
}

const (
	QUEUE_TYPE_AT_LEAST_ONCE = iota
	QUEUE_TYPE_AT_MOST_ONCE
)

type Options struct {
	Url       string
	ClusterId string
}

type Queue interface {
	Open(ctx context.Context, clientId string, options *Options) error
	Close() error
	Publish(ctx context.Context, qt QueueType, topic string, msg *Message) error
	Subscribe(qt QueueType, topic string, receiverChan chan<- []byte) error
}
