package stan

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
)

const (
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

	log.DbgF("[STAN] connected client_id %s", clientId)

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

func (s *stanImpl) Publish(ctx context.Context, topic string, msg *queue.Message) error {
	if s.conn == nil {
		return errors.New("trying to publish to undefined connection")
	}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.TrcF("[STAN] published: %s", string(m))
	return s.conn.Publish(topic, m)
}

func (s *stanImpl) Subscribe(topic string, receiverChan chan<- []byte) error {
	_, err := s.conn.Subscribe(topic, func(m *stan.Msg) {
		log.DbgF("[STAN] received: %s", string(m.Data))
		receiverChan <- m.Data
	}, stan.DurableName(s.clientId))
	return err
}

func (s *stanImpl) PublishAtMostOnce(ctx context.Context, topic string, msg *queue.Message) error {
	if s.conn == nil {
		return errors.New("trying to publish to undefined connection")
	}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.TrcF("[NATS] published: %s", string(m))
	return s.conn.NatsConn().Publish(topic, m)
}

func (s *stanImpl) SubscribeAtMostOnce(topic string, receiverChan chan<- []byte) error {
	_, err := s.conn.NatsConn().Subscribe(topic, func(m *nats.Msg) {
		log.TrcF("[NATS] received: %s", string(m.Data))
		receiverChan <- m.Data
	})
	return err
}