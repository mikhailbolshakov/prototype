package main

import (
	"gitlab.medzdrav.ru/prototype/api"
	"gitlab.medzdrav.ru/prototype/tasks"
	"gitlab.medzdrav.ru/prototype/users"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		usersService = users.NewService()
		tasksService = tasks.NewService()
		apiService   = api.NewService()
	)

	err := usersService.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = tasksService.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = apiService.Start()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

}