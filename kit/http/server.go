package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Srv          *http.Server
	AuthRouter   *mux.Router
	NoAuthRouter *mux.Router
}

type RouteSetter interface {
	Set(authRouter, noAuthRouter *mux.Router)
}

func NewHttpServer(host, port string, routeSetters ...RouteSetter) *Server {

	r := mux.NewRouter()
	noAuthRouter := r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()
	authRouter := r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	for _, rs := range routeSetters {
		rs.Set(authRouter, noAuthRouter)
	}

	s := &Server{
		Srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      r,
			WriteTimeout: time.Hour,
			ReadTimeout:  time.Hour,
		},
		AuthRouter: authRouter,
		NoAuthRouter: noAuthRouter,
	}

	return s
}

func (s *Server) SetAuthMiddleware(mdws ...mux.MiddlewareFunc) {

	for _, m := range mdws {
		s.AuthRouter.Use(m)
	}

}

func (s *Server) Open() error {
	return s.Srv.ListenAndServe()
}
