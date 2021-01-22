package services

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

	authRouter.HandleFunc("/api/{userId}/services/balance", func(writer http.ResponseWriter, request *http.Request) {
		c.AddBalance(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/{userId}/services/balance", func(writer http.ResponseWriter, request *http.Request) {
		c.GetBalance(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/{userId}/services/delivery", func(writer http.ResponseWriter, request *http.Request) {
		c.Delivery(writer, request)
	}).Methods("POST")


}
