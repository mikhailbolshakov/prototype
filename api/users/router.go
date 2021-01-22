package users

import (
	"github.com/gorilla/mux"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	"log"
	"net/http"
)

type Router struct {
	auth auth.AuthenticationHandler
}

func New(auth auth.AuthenticationHandler) http2.RouteSetter {
	return &Router{
		auth: auth,
	}
}

func (u *Router) Set(authRouter, noAuthRouter *mux.Router) {

	c, err := newController(u.auth)
	if err != nil {
		log.Fatalln(err)
		return
	}

	authRouter.HandleFunc("/api/users", func(writer http.ResponseWriter, request *http.Request) {
		c.Create(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/users/username/{un}", func(writer http.ResponseWriter, request *http.Request) {
		c.GetByUsername(writer, request)
	}).Methods("GET")

	authRouter.HandleFunc("/api/users", func(writer http.ResponseWriter, request *http.Request) {
		c.Search(writer, request)
	}).Methods("GET")

	noAuthRouter.HandleFunc("/api/users/login", func(writer http.ResponseWriter, request *http.Request) {
		c.Login(writer, request)
	}).Methods("POST")

}
