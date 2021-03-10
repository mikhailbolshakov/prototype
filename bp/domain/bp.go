package domain

import "gitlab.medzdrav.ru/prototype/kit/queue/listener"

type BusinessProcess interface {
	Init() error
	GetId() string
	// GetBPMNFileName returns bpmn file name relative to {bpmn.src-folder} config entry
	GetBPMNFileName() string
	SetQueueListeners(ql listener.QueueListener)
}
