package services

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

	r.HandleFunc("/api/{userId}/services", func(writer http.ResponseWriter, request *http.Request) {
		c.AddUserServices(writer, request)
	}).Methods("POST")

	r.HandleFunc("/api/{userId}/services/writeoff", func(writer http.ResponseWriter, request *http.Request) {
		c.WriteOff(writer, request)
	}).Methods("POST")

	r.HandleFunc("/api/{userId}/services/lock", func(writer http.ResponseWriter, request *http.Request) {
		c.Lock(writer, request)
	}).Methods("POST")

	r.HandleFunc("/api/{userId}/services/balance", func(writer http.ResponseWriter, request *http.Request) {
		c.GetBalance(writer, request)
	}).Methods("GET")

	r.HandleFunc("/api/{userId}/services/delivery", func(writer http.ResponseWriter, request *http.Request) {
		c.Delivery(writer, request)
	}).Methods("POST")


}
