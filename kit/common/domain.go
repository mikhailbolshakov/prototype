package common

import (
	"context"
	"fmt"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

type BaseService struct {
	Queue queue.Queue
}

func (s *BaseService) Publish(ctx context.Context, o interface{}, qt queue.QueueType, topic string) error {

	m := &queue.Message{ Payload: o	}

	if rCtx, ok := kitContext.Request(ctx); ok {
		m.Ctx = rCtx
	} else {
		return fmt.Errorf("cannot publish to queue topic %s, context invalid", topic)
	}

	return s.Queue.Publish(ctx, qt, topic, m)

}
