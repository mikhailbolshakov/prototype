package logger

import (
	"gitlab.medzdrav.ru/prototype/api/meta"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

var Logger = log.Init(log.TraceLevel)

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv(meta.Meta.InstanceId())
	}
}

func L() log.CLogger {
	return LF()()
}
