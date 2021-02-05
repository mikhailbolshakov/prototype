package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var logger = logrus.New()

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = "panic"
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = "fatal"
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = "error"
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = "warning"
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = "info"
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = "debug"
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = "trace"
)

func Init(level string) error {

	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(lv)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	//logger.SetReportCaller(true)

	return nil
}

func Inf(args ...interface{}) {
	logger.Infoln(args...)
}

func InfF(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Dbg(args ...interface{}) {
	logger.Debugln(args...)
}

func DbgF(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Trc(args ...interface{}) {
	logger.Traceln(args...)
}

func TrcF(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Err(err error, stack bool) {

	if stack {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)
		logger.Errorln(err.Error(), fmt.Sprintf("%s", buf))
	} else {
		logger.Errorln(err.Error())
	}

}

