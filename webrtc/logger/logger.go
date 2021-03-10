package logger

import (
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
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
