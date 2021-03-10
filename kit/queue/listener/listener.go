package listener

import (
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"sync"
)

type QueueMessageHandler func(payload []byte) error

type QueueListener interface {
	Add(qt queue.QueueType, topic string, h ...QueueMessageHandler)
	ListenAsync()
	Stop()
	Clear()
}

func NewQueueListener(q queue.Queue, logger log.CLoggerFunc) QueueListener {

	th := map[queue.QueueType]map[string][]QueueMessageHandler{}
	th[queue.QUEUE_TYPE_AT_LEAST_ONCE] = make(map[string][]QueueMessageHandler)
	th[queue.QUEUE_TYPE_AT_MOST_ONCE] = make(map[string][]QueueMessageHandler)

	return &queueListener{
		topicHandlers: th,
		listening:     false,
		queue:         q,
		logger:        logger,
	}
}

type queueListener struct {
	sync.RWMutex
	queue         queue.Queue
	topicHandlers map[queue.QueueType]map[string][]QueueMessageHandler
	quit          chan struct{}
	listening     bool
	logger        log.CLoggerFunc
}

func (q *queueListener) Add(qt queue.QueueType, topic string, h ...QueueMessageHandler) {

	q.Stop()

	q.Lock()
	defer q.Unlock()

	var handlers []QueueMessageHandler
	handlers, ok := q.topicHandlers[qt][topic]
	if !ok {
		handlers = []QueueMessageHandler{}
	}

	for _, hnd := range h {
		handlers = append(handlers, hnd)
	}
	q.topicHandlers[qt][topic] = handlers

}

func (q *queueListener) ListenAsync() {

	for queueType, topicHandlers := range q.topicHandlers {
		for topic, handlers := range topicHandlers {
			go func(qt queue.QueueType, tp string, hnds []QueueMessageHandler) {
				c := make(chan []byte)
				_ = q.queue.Subscribe(qt, tp, c)
				for {
					select {
					case msg := <-c:
						for _, h := range hnds {
							m := msg
							h := h
							go func() {
								l := q.logger().Pr("queue").Cmp("listener").F(log.FF{"topic": tp})
								l.TrcF("%s", string(m))
								if err := h(m); err != nil {
									l.E(err).St().Err()
								}
							}()
						}
					case <-q.quit:
						return
					}
				}
			}(queueType, topic, handlers)
		}
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
	q.topicHandlers[queue.QUEUE_TYPE_AT_LEAST_ONCE] = make(map[string][]QueueMessageHandler)
	q.topicHandlers[queue.QUEUE_TYPE_AT_MOST_ONCE] = make(map[string][]QueueMessageHandler)
}
