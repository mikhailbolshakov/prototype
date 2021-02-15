package chat

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

	authRouter.HandleFunc("/api/chat/users/{id}/status/{status}", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.SetStatus(writer, request)
	}).Methods("PUT")

	authRouter.HandleFunc("/api/chat/users/{id}/posts", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.Post(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/chat/users/{id}/posts/ephemeral", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.EphemeralPost(writer, request)
	}).Methods("POST")

}
