package stan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

const (
	// TODO: cfg
	CLUSTER_ID = "test-cluster"
)

type stanImpl struct {
	conn stan.Conn
	clientId string
}

func New() queue.Queue {
	return &stanImpl{}
}

func (s *stanImpl) Open(ctx context.Context, clientId string) error {

	s.clientId = clientId
	c, err := stan.Connect(CLUSTER_ID, clientId)
	if err != nil {
		return err
	}
	s.conn = c

	log.DbgF("[nats] connected client_id %s\n", clientId)

	return nil
}

func (s *stanImpl) Close() error {
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		return err
	}
	return nil
}

func (s *stanImpl) Publish(ctx context.Context, qt queue.QueueType, topic string, msg *queue.Message) error {
	if s.conn == nil {
		return errors.New("trying to publish to undefined connection")
	}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.TrcF("[nats] published: %s\n", string(m))

	if qt == queue.QUEUE_TYPE_AT_LEAST_ONCE {
		return s.conn.Publish(topic, m)
	} else if qt == queue.QUEUE_TYPE_AT_MOST_ONCE {
		return s.conn.NatsConn().Publish(topic, m)
	} else {
		return fmt.Errorf("[nats] not supported queue type")
	}

}

func (s *stanImpl) Subscribe(qt queue.QueueType, topic string, receiverChan chan<- []byte) error {

	if qt == queue.QUEUE_TYPE_AT_LEAST_ONCE {
		_, err := s.conn.Subscribe(topic, func(m *stan.Msg) {
			log.DbgF("[nats][at-least-once] received: %s\n", string(m.Data))
			receiverChan <- m.Data
		}, stan.DurableName(s.clientId))
		return err
	} else if  qt == queue.QUEUE_TYPE_AT_MOST_ONCE {
		_, err := s.conn.NatsConn().Subscribe(topic, func(m *nats.Msg) {
			log.TrcF("[nats][at-most-once] received: %s\n", string(m.Data))
			receiverChan <- m.Data
		})
		return err
	} else {
		return fmt.Errorf("[nats] not supported queue type")
	}

}
