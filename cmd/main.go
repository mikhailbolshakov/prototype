package main

import (
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
)

func main() {

	logger := log.Init(log.TraceLevel)

	// get empty context
	// TODO: change to some custom background context
	ctx := kitContext.NewRequestCtx().Empty().ToContext(nil)

	// instantiate all services
	srvs := []service.Service{
		config.New(),
		users.New(),
		tasks.New(),
		chat.New(),
		services.New(),
		bp.New(),
		api.New(),
		webrtc.New(),
		sessions.New(),
	}

	for _, s := range srvs{
		s := s
		go func() {
			l := log.L(logger).Srv(s.GetCode()).Cmp("main")
			err := s.Init(ctx)
			if err != nil {
				l.E(err).St().Err()
				os.Exit(1)
			}
			if err := s.ListenAsync(ctx); err != nil {
				l.E(err).St().Err()
				os.Exit(1)
			}
		}()
	}

	// handle app close
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.L(logger).Inf("quit signal")
	for _, s := range srvs {
		s.Close(ctx)
	}
	os.Exit(0)

}