package tasks

import (
	"github.com/gorilla/mux"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"log"
	"net/http"
)

type Router struct {}

func New() http2.RouteSetter {
	return &Router{}
}

func (u *Router) Set(authRouter, noAuthRouter *mux.Router) {

	c, err := newController()
	if err != nil {
		log.Fatalln(err)
		return
	}

	authRouter.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		c.New(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks/{id}/transitions/{transitionId}", func(writer http.ResponseWriter, request *http.Request) {
		c.MakeTransition(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		c.GetById(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/{id}/assignee", func(writer http.ResponseWriter, request *http.Request) {
		c.SetAssignee(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		c.Search(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/assignment/log", func(writer http.ResponseWriter, request *http.Request) {
		c.AssignmentLog(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/{id}/history", func(writer http.ResponseWriter, request *http.Request) {
		c.GetHistory(writer, request)
	}).Methods("GET")

}
