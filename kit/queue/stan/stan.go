package stan

import (
	"errors"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	CLUSTER_ID = "test-cluster"
)

type Stan struct {
	conn stan.Conn
	clientId string
}

func (s *Stan) Open(clientId string) error {

	s.clientId = clientId
	c, err := stan.Connect(CLUSTER_ID, clientId)
	if err != nil {
		return err
	}
	s.conn = c

	log.Printf("[STAN] connected client_id %s", clientId)

	return nil
}

func (s *Stan) Close() error {
	if s.conn != nil {
		s.conn = nil
		return s.Close()
	}
	return nil
}

func (s *Stan) Publish(topic string, msg []byte) error {
	if s.conn == nil {
		return errors.New("trying to publish to undefined connection")
	}
	log.Printf("[STAN] published: %s", string(msg))
	return s.conn.Publish(topic, msg)
}

func (s *Stan) Subscribe(topic string, receiverChan chan<- []byte) error {
	_, err := s.conn.Subscribe(topic, func(m *stan.Msg) {
		log.Printf("[STAN] received: %s", string(m.Data))
		receiverChan <- m.Data
	}, stan.DurableName(s.clientId))
	return err
}