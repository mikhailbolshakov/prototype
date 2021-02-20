package log

import (
	"fmt"
	context2 "gitlab.medzdrav.ru/prototype/kit/context"
	"testing"
)


func Test_Clogger_WithCtx(t *testing.T) {
	_ = Init(TraceLevel)
	ctx := context2.NewRequestCtx().Rest().WithNewRequestId().WithUser("1", "john").ToContext(nil)
	l := L().C(ctx)
	l.Inf("I'm logger")
}

func Test_Clogger_WithComponentAndMethod(t *testing.T) {
	_ = Init(TraceLevel)
	l := L().Cmp("service").Mth("do")
	l.Inf("I'm logger")
}

func Test_Clogger_WithComponentMethodAndCtx(t *testing.T) {
	_ = Init(TraceLevel)
	ctx := context2.NewRequestCtx().Rest().WithNewRequestId().WithUser("1", "john").ToContext(nil)
	l := L().Cmp("service").Mth("do").C(ctx)
	l.Inf("I'm logger")
}

func Test_Clogger_All(t *testing.T) {
	_ = Init(TraceLevel)
	ctx := context2.NewRequestCtx().Rest().WithNewRequestId().WithUser("1", "john").ToContext(nil)
	l := L().Cmp("service").Mth("do").C(ctx).F(FF{"field": "value"})
	l.Inf("I'm logger")
}

func Test_Clogger_WithFields(t *testing.T) {
	_ = Init(TraceLevel)
	l := L().F(FF{"field": "value"})
	l.Inf("I'm logger")
}

func Test_Clogger_WithErr(t *testing.T) {
	_ = Init(TraceLevel)
	l := L().E(fmt.Errorf("error"))
	l.Err("my bad")
}

func Test_Clogger_WithErrStack(t *testing.T) {
	_ = Init(TraceLevel)
	l := L().E(fmt.Errorf("error")).St()
	l.Err("my bad")
}