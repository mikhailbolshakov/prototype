package users

import (
	"github.com/gorilla/mux"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Router struct {
	ctrl Controller
}

func NewRouter(c Controller) http2.RouteSetter {
	return &Router{
		ctrl: c,
	}
}

func (r *Router) Set(authRouter, noAuthRouter *mux.Router) {

	authRouter.HandleFunc("/api/users/client", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.CreateClient(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/users/consultant", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.CreateConsultant(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/users/expert", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.CreateExpert(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/users/username/{un}", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.GetByUsername(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/users", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.Search(writer, request)
	}).Methods("GET")


}
