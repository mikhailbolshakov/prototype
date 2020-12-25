package users

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

	r.HandleFunc("/api/users", func(writer http.ResponseWriter, request *http.Request) {
		c.Create(writer, request)
	}).Methods("POST")

	r.HandleFunc("/api/users/username/{un}", func(writer http.ResponseWriter, request *http.Request) {
		c.GetByUsername(writer, request)
	}).Methods("GET")

}
