package bpm

type Params struct {
	Port string
	Host string
}

type Engine interface {
	Open(params *Params) error
	Close() error
	IsOpened() bool
	DeployBPMNs(paths []string) error
	RegisterTaskHandlers(handlers map[string]interface{}) error
	StartProcess(processId string, vars map[string]interface{}) (string, error)
	SendMessage(messageId string, correlationId string, vars map[string]interface{}) error
	SendError(jobId int64, errCode, errMessage string) error
}
