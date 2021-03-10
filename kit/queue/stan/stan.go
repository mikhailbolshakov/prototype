package stan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

type stanImpl struct {
	conn     stan.Conn
	clientId string
	logger   log.CLoggerFunc
}

func New(logger log.CLoggerFunc) queue.Queue {
	return &stanImpl{
		logger: logger,
	}
}

func (s *stanImpl) l() log.CLogger {
	return s.logger().Pr("queue").Cmp("stan")
}

func (s *stanImpl) Open(ctx context.Context, clientId string, options *queue.Options) error {

	l := s.l().Mth("open").F(log.FF{"client": clientId, "url": options.Url}).Dbg("connecting")

	s.clientId = clientId
	c, err := stan.Connect(options.ClusterId, clientId, stan.NatsURL(options.Url))
	if err != nil {
		return err
	}
	s.conn = c

	l.Inf("ok")

	return nil
}

func (s *stanImpl) Close() error {
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		if err != nil {
			return err
		}
		s.l().Mth("close").Inf("closed")
	}
	return nil
}

func (s *stanImpl) Publish(ctx context.Context, qt queue.QueueType, topic string, msg *queue.Message) error {

	l := s.l().Mth("publish").F(log.FF{"topic": topic, "type": qt.String()})

	if msg.Ctx == nil {
		msg.Ctx = kitContext.NewRequestCtx().Queue().WithNewRequestId()
	}
	l.C(msg.Ctx.ToContext(context.Background()))

	if s.conn == nil {
		return fmt.Errorf("no open connection")
	}

	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	l.Dbg("ok")
	l.TrcF("%s\n", string(m))

	if qt == queue.QUEUE_TYPE_AT_LEAST_ONCE {
		return s.conn.Publish(topic, m)
	} else if qt == queue.QUEUE_TYPE_AT_MOST_ONCE {
		return s.conn.NatsConn().Publish(topic, m)
	} else {
		return fmt.Errorf("queue type %d not supported", qt)
	}

}

func (s *stanImpl) Subscribe(qt queue.QueueType, topic string, receiverChan chan<- []byte) error {

	l := s.l().Mth("received").F(log.FF{"topic": topic, "type": qt.String()})

	if qt == queue.QUEUE_TYPE_AT_LEAST_ONCE {
		_, err := s.conn.Subscribe(topic, func(m *stan.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		}, stan.DurableName(s.clientId))
		return err
	} else if qt == queue.QUEUE_TYPE_AT_MOST_ONCE {
		_, err := s.conn.NatsConn().Subscribe(topic, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		})
		return err
	} else {
		return fmt.Errorf("queue type %d not supported", qt)
	}

}
