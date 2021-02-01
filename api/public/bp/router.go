package bp

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

	authRouter.HandleFunc("/api/bp/start", func(writer http.ResponseWriter, request *http.Request) {
		r.c.StartProcess(writer, request)
	}).Methods("POST")

}
