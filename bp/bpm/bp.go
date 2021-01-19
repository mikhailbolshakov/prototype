package bpm

import "gitlab.medzdrav.ru/prototype/kit/queue/listener"

type BusinessProcess interface {
	Init() error
	GetId() string
	SetQueueListeners(ql listener.QueueListener)
}
