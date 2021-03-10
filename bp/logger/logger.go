package logger

import (
	"gitlab.medzdrav.ru/prototype/bp/meta"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

var Logger = log.Init(log.TraceLevel)

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv(meta.ServiceCode).Nd(meta.NodeId)
	}
}

func L() log.CLogger {
	return LF()()
}
