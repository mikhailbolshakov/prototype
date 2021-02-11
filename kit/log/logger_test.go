package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_ErrorWithStack(t *testing.T) {
	Err(fmt.Errorf("error"), true)
}

func Test_WithFields(t *testing.T) {
	l := GetLogger()
	loggerCtx := l.WithFields(logrus.Fields{"traceId": "123", "host": "127.0.0.1"})
	loggerCtx.Info("message 1")
	loggerCtx.Info("message 2")
}

func Test_WithContext(t *testing.T) {
	l := GetLogger()
	ctx := context.WithValue(context.Background(), "TraceId", "123")
	le := l.WithContext(ctx)
	le.Info("message")
}