package mm

import (
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
)

func NewQueueHandler() queue.Queue {
	return &stan.Stan{}
}
