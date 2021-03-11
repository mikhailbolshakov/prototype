package log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"os"
	"runtime"
)

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

type Logger struct {
	Logrus *logrus.Logger
}

func Init(level string) *Logger {

	logger := &Logger{
		Logrus: logrus.New(),
	}
	logger.SetLevel(level)
	logger.Logrus.SetOutput(os.Stdout)
	logger.Logrus.SetFormatter(&Formatter{
		FixedFields:            []string{"service", "node", "protocol", "component", "method"},
		TimestampFormat:        "2006-01-02T15:04:05-0700",
		HideKeysForFixedFields: true,
		NoColors:               true,
		NoFieldsColors:         true,
		NoFieldsSpace:          true,
	})

	return logger
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.Logrus
}

func (l *Logger) SetLevel(level string) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	l.Logrus.SetLevel(lv)
}

type CLoggerFunc func() CLogger

// CLogger provides structured logging abilities
// !!!! Not thread safe. Don't share one CLogger instance through multiple goroutines
type CLogger interface {
	C(ctx context.Context) CLogger // C - adds request context to log
	F(fields FF) CLogger           // F - adds fields to log
	E(err error) CLogger           // E - adds error to log
	St() CLogger                   // St - adds stack to log (if err is already set)
	Cmp(c string) CLogger          // Cmp - adds component
	Mth(m string) CLogger          // Mth - adds method
	Pr(m string) CLogger           // Pr - adds protocol
	Srv(s string) CLogger          // Srv - adds service code
	Nd(n string) CLogger           // Nd - adds node code
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
	Clone() CLogger
	Printf(string, ...interface{})
}

func L(logger *Logger) CLogger {
	return &clogger{
		logger: logger,
		lre: logrus.NewEntry(logger.Logrus),
	}
}

type clogger struct {
	logger *Logger
	lre *logrus.Entry
	err error
}

// always use Clone when pass CLogger between goroutines
func (cl *clogger) Clone() CLogger {
	entry := logrus.NewEntry(cl.logger.Logrus)
	if len(cl.lre.Data) > 0 {
		marshaled, _ := json.Marshal(cl.lre.Data)
		_ = json.Unmarshal(marshaled, &entry.Data)
	}
	clone := &clogger{
		lre: entry,
		err: cl.err,
	}
	return clone
}

func (cl *clogger) C(ctx context.Context) CLogger {
	if r, ok := kitContext.Request(ctx); ok {
		cl.F(FF{"ctx-cl": r.GetClientType(),
				"ctx-rid": r.GetRequestId(),
				"ctx-un": r.GetUsername(),
				"ctx-sid": r.GetSessionId(),
				})
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

func (cl *clogger) Srv(s string) CLogger {
	cl.lre = cl.lre.WithField("service", s)
	return cl
}

func (cl *clogger) Nd(s string) CLogger {
	cl.lre = cl.lre.WithField("node", s)
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
	cl.lre.Tracef(format, args...)
	return cl
}

func (cl *clogger) Fatal(args ...interface{}) CLogger {
	cl.lre.Fatalln(args...)
	return cl
}

func (cl *clogger) FatalF(format string, args ...interface{}) CLogger {
	cl.lre.Fatalf(format, args...)
	return cl
}

func (cl *clogger) Printf(f string, args ...interface{}) {
	cl.DbgF(f, args...)
}
