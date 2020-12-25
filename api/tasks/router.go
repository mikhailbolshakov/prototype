package tasks

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Router struct {}

func (u *Router) Set(r *mux.Router) {

	c, err := newController()
	if err != nil {
		log.Fatalln(err)
		return
	}

	r.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		c.New(writer, request)
	}).Methods("POST")

	r.HandleFunc("/api/tasks/{id}/transitions/{transitionId}", func(writer http.ResponseWriter, request *http.Request) {
		c.MakeTransition(writer, request)
	}).Methods("POST")

}
