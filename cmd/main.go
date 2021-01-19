package main

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/api"
	"gitlab.medzdrav.ru/prototype/bp"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/mm"
	"gitlab.medzdrav.ru/prototype/services"
	"gitlab.medzdrav.ru/prototype/tasks"
	"gitlab.medzdrav.ru/prototype/users"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	srvs := []service.Service{
		users.New(),
		tasks.New(),
		mm.New(),
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

	// run listeners
	for _, s := range srvs {
		if err := s.ListenAsync(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}