package common

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"log"
)

type BaseService struct {
	Queue queue.Queue
}

func (s *BaseService) Publish(o interface{}, topic string) {
	go func() {
		j, err := json.Marshal(o)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = s.Queue.Publish(topic, j)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
}
