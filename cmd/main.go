package main

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/api"
	"gitlab.medzdrav.ru/prototype/bp"
	"gitlab.medzdrav.ru/prototype/chat"
	"gitlab.medzdrav.ru/prototype/config"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/services"
	"gitlab.medzdrav.ru/prototype/sessions"
	"gitlab.medzdrav.ru/prototype/tasks"
	"gitlab.medzdrav.ru/prototype/users"
	"gitlab.medzdrav.ru/prototype/webrtc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	if err := log.Init(log.TraceLevel); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// get empty context
	// TODO: change to some custom background context
	ctx := kitContext.NewRequestCtx().Empty().ToContext(nil)

	l := log.L().Cmp("main")

	// load config first
	cfg := config.New()
	if err := cfg.Init(ctx); err != nil {
		l.E(err).St().Err()
		os.Exit(1)
	}
	if err := cfg.ListenAsync(ctx); err != nil {
		l.E(err).St().Err()
		os.Exit(1)
	}

	// to avoid the case when services start earlier than config server
	// we need a retry approach to avoid this for microservices
	time.Sleep(time.Second)

	// instantiate all services
	srvs := []service.Service{
		users.New(),
		tasks.New(),
		chat.New(),
		services.New(),
		bp.New(),
		api.New(),
		webrtc.New(),
		sessions.New(),
	}

	// init service
	for _, s := range srvs {
		err := s.Init(ctx)
		if err != nil {
			l.E(err).St().Err()
			os.Exit(1)
		}
	}

	//run listeners
	for _, s := range srvs {
		if err := s.ListenAsync(ctx); err != nil {
			l.E(err).St().Err()
			os.Exit(1)
		}
	}

	// handle app close
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Inf("quit signal")
	for _, s := range srvs {
		s.Close(ctx)
	}
	cfg.Close(ctx)
	os.Exit(0)

}