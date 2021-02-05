package main

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/api"
	"gitlab.medzdrav.ru/prototype/bp"
	"gitlab.medzdrav.ru/prototype/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/chat"
	"gitlab.medzdrav.ru/prototype/services"
	"gitlab.medzdrav.ru/prototype/tasks"
	"gitlab.medzdrav.ru/prototype/users"
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

	// load config first
	cfg := config.New()
	if err := cfg.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := cfg.ListenAsync(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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
	}

	// init service
	for _, s := range srvs {
		err := s.Init()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	//run listeners
	for _, s := range srvs {
		if err := s.ListenAsync(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	// handle app close
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	for _, s := range srvs {
		s.Close()
	}
	cfg.Close()
	os.Exit(0)

}