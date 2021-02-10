package domain

import "gitlab.medzdrav.ru/prototype/kit/queue/listener"

type BusinessProcess interface {
	Init() error
	GetId() string
	GetBPMNPath() string
	SetQueueListeners(ql listener.QueueListener)
}
