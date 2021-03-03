package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
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

func GetLogger() *logrus.Logger {
	return logger
}

func Init(level string) error {

	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(lv)
	logger.SetFormatter(&Formatter{
		FixedFields:            []string{"protocol", "component", "method"},
		TimestampFormat:        "2006-01-02T15:04:05-0700",
		HideKeysForFixedFields: true,
		NoColors:               true,
		NoFieldsColors:         true,
		NoFieldsSpace:          true,
	})
	return nil
}

// CLogger provides rich logging abilities
// !!!! Not thread safe. Don't share one CLogger instance through multiple goroutines
type CLogger interface {
	C(ctx context.Context) CLogger // C - adds request context to log
	F(fields FF) CLogger           // F - adds fields to log
	E(err error) CLogger           // E - adds error to log
	St() CLogger                   // St - adds stack to log (if err is already set)
	Cmp(c string) CLogger          // Cmp - adds component
	Mth(m string) CLogger          // Mth - adds method
	Pr(m string) CLogger           // Pr - adds protocol
	Inf(args ...interface{}) CLogger
	InfF(format string, args ...interface{}) CLogger
	Err(args ...interface{}) CLogger
	ErrF(format string, args ...interface{}) CLogger
	Dbg(args ...interface{}) CLogger
	DbgF(format string, args ...interface{}) CLogger
	Trc(args ...interface{}) CLogger
	TrcF(format string, args ...interface{}) CLogger
	Warn(args ...interface{}) CLogger
	WarnF(format string, args ...interface{}) CLogger
	Fatal(args ...interface{}) CLogger
	FatalF(format string, args ...interface{}) CLogger
}

func L() CLogger {
	return &clogger{
		lre: logrus.NewEntry(logger),
	}
}

type clogger struct {
	lre  *logrus.Entry
	err  error
	tags []string
}

func (cl *clogger) C(ctx context.Context) CLogger {
	if r, ok := kitContext.Request(ctx); ok {
		cl.F(FF{"rid": r.GetRequestId(), "un": r.GetUsername(), "cl": r.GetClientType()})
	}
	return cl
}

type FF map[string]interface{}

func (cl *clogger) F(fields FF) CLogger {
	cl.lre = cl.lre.WithFields(map[string]interface{}(fields))
	return cl
}

func (cl *clogger) E(err error) CLogger {
	cl.lre = cl.lre.WithError(err)
	cl.err = err
	return cl
}

func (cl *clogger) St() CLogger {
	if cl.err != nil {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, false)
		cl.lre = cl.lre.WithField("st", fmt.Sprintf("%s", buf))
	}
	return cl
}

func (cl *clogger) Cmp(c string) CLogger {
	cl.lre = cl.lre.WithField("component", c)
	return cl
}

func (cl *clogger) Pr(c string) CLogger {
	cl.lre = cl.lre.WithField("protocol", c)
	return cl
}

func (cl *clogger) Mth(m string) CLogger {
	cl.lre = cl.lre.WithField("method", m)
	return cl
}

func (cl *clogger) Err(args ...interface{}) CLogger {
	cl.lre.Errorln(args...)
	return cl
}

func (cl *clogger) ErrF(format string, args ...interface{}) CLogger {
	cl.lre.Errorf(format, args...)
	return cl
}

func (cl *clogger) Inf(args ...interface{}) CLogger {
	cl.lre.Infoln(args...)
	return cl
}

func (cl *clogger) InfF(format string, args ...interface{}) CLogger {
	cl.lre.Infof(format, args...)
	return cl
}

func (cl *clogger) Warn(args ...interface{}) CLogger {
	cl.lre.Warningln(args...)
	return cl
}

func (cl *clogger) WarnF(format string, args ...interface{}) CLogger {
	cl.lre.Warningf(format, args...)
	return cl
}

func (cl *clogger) Dbg(args ...interface{}) CLogger {
	cl.lre.Debugln(args...)
	return cl
}

func (cl *clogger) DbgF(format string, args ...interface{}) CLogger {
	cl.lre.Debugf(format, args...)
	return cl
}

func (cl *clogger) Trc(args ...interface{}) CLogger {
	cl.lre.Traceln(args...)
	return cl
}

func (cl *clogger) TrcF(format string, args ...interface{}) CLogger {
	cl.lre.Tracef(format, args)
	return cl
}

func (cl *clogger) Fatal(args ...interface{}) CLogger {
	cl.lre.Fatalln(args...)
	return cl
}

func (cl *clogger) FatalF(format string, args ...interface{}) CLogger{
	cl.lre.Fatalf(format, args...)
	return cl
}
