package webrtc

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

	authRouter.HandleFunc("/api/rooms", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.CreateRoom(writer, request)
	}).Methods("POST")

	authRouter.HandleFunc("/api/rooms/{roomId}", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.GetRoom(writer, request)
	}).Methods("GET")

}
