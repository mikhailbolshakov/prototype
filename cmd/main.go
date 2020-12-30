package main

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/api"
	"gitlab.medzdrav.ru/prototype/kit/service"
	"gitlab.medzdrav.ru/prototype/tasks"
	"gitlab.medzdrav.ru/prototype/users"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	services := []service.Service{
		users.New(),
		tasks.New(),
		api.New(),
	}

	// init service
	for _, s := range services {
		err := s.Init()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	// run listeners
	for _, s := range services {
		if err := s.ListenAsync(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}