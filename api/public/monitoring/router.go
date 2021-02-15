package monitoring

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

	noAuthRouter.HandleFunc("/api/monitor/sessions/users/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.GetUserSessions(writer, request)
	}).Methods("GET")

	noAuthRouter.HandleFunc("/api/monitor/sessions", func(writer http.ResponseWriter, request *http.Request) {
		r.ctrl.GetTotalSessions(writer, request)
	}).Methods("GET")

}
