package main

import (
	"gitlab.medzdrav.ru/prototype/bp/logger"
	"gitlab.medzdrav.ru/prototype/bp"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// init context
	ctx := kitContext.NewRequestCtx().Empty().WithNewRequestId().ToContext(nil)

	// create a new service
	s := bp.New()

	l := logger.L().Mth("main").Inf("created")

	// init service
	if err := s.Init(ctx); err != nil {
		l.E(err).St().Err("initialization")
		os.Exit(1)
	}

	l.Inf("initialized")

	// start listening
	if err := s.ListenAsync(ctx); err != nil {
		l.E(err).St().Err("listen")
		os.Exit(1)
	}

	l.Inf("listening")

	// handle app close
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Inf("graceful shutdown")
	s.Close(ctx)
	os.Exit(0)
}
