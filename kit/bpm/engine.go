package bpm

type Engine interface {
	Open() error
	Close() error
	IsOpened() bool
	DeployBPMNs(paths []string) error
	RegisterTaskHandlers(handlers map[string]interface{}) error
	StartProcess(processId string, vars map[string]interface{}) (string, error)
	SendMessage(messageId string, correlationId string, vars map[string]interface{}) error
}
