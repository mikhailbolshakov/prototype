package listener

import (
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"log"
	"sync"
)

type QueueMessageHandler func(payload []byte) error

type QueueListener interface {
	Add(topic string, h ...QueueMessageHandler)
	ListenAsync()
	Stop()
	Clear()
}

func NewQueueListener(queue queue.Queue) QueueListener {
	return &queueListener{
		topicHandlers: make(map[string][]QueueMessageHandler),
		listening:     false,
		queue:         queue,
	}
}

type queueListener struct {
	sync.RWMutex
	queue         queue.Queue
	topicHandlers map[string][]QueueMessageHandler
	quit          chan struct{}
	listening     bool
}

func (q *queueListener) Add(topic string, h ...QueueMessageHandler) {

	q.Stop()

	q.Lock()
	defer q.Unlock()

	var handlers []QueueMessageHandler
	handlers, ok := q.topicHandlers[topic]
	if !ok {
		handlers = []QueueMessageHandler{}
	}

	for _, hnd := range h {
		handlers = append(handlers, hnd)
	}
	q.topicHandlers[topic] = handlers

}

func (q *queueListener) ListenAsync() {

	for topic, handlers := range q.topicHandlers {
		handlers := handlers
		go func(t string, hnds []QueueMessageHandler) {
			c := make(chan []byte)
			_ = q.queue.Subscribe(t, c)
			for {
				select {
				case msg := <-c:
					for _, h := range hnds {
						m := msg
						h := h
						go func() {
							if err := h(m); err != nil {
								log.Printf("[ERROR] %v handler error %v", t, err)
							}
						}()
					}
				case <-q.quit:
					return
				}
			}
		}(topic, handlers)

	}

}

func (q *queueListener) Stop() {
	q.RLock()
	l := q.listening
	q.RUnlock()

	if l {
		q.quit <- struct{}{}
		q.Lock()
		defer q.Unlock()
		q.listening = false
	}
}

func (q *queueListener) Clear() {
	q.Stop()
	q.quit <- struct{}{}
	q.Lock()
	defer q.Unlock()
	q.topicHandlers = make(map[string][]QueueMessageHandler)
}
