package tasks

import (
	"github.com/gorilla/mux"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Router struct {
	c Controller
}

func NewRouter(c Controller) http2.RouteSetter {
	return &Router{
		c: c,
	}
}

func (r *Router) Set(authRouter, noAuthRouter *mux.Router) {

	authRouter.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		r.c.New(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks/{id}/transitions/{transitionId}", func(writer http.ResponseWriter, request *http.Request) {
		r.c.MakeTransition(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r.c.GetById(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/{id}/assignee", func(writer http.ResponseWriter, request *http.Request) {
		r.c.SetAssignee(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/tasks", func(writer http.ResponseWriter, request *http.Request) {
		r.c.Search(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/assignment/log", func(writer http.ResponseWriter, request *http.Request) {
		r.c.AssignmentLog(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/tasks/{id}/history", func(writer http.ResponseWriter, request *http.Request) {
		r.c.GetHistory(writer, request)
	}).Methods("GET")

}
