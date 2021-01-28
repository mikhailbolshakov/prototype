package services

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

	authRouter.HandleFunc("/api/{userId}/services/balance", func(writer http.ResponseWriter, request *http.Request) {
		r.c.AddBalance(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/{userId}/services/balance", func(writer http.ResponseWriter, request *http.Request) {
		r.c.GetBalance(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/{userId}/services/delivery", func(writer http.ResponseWriter, request *http.Request) {
		r.c.Delivery(writer, request)
	}).Methods("POST")


}
