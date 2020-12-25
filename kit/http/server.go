package http

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"time"
)

type Server struct {
	Srv *http.Server
}

type RouteSetter interface {
	Set(r *mux.Router)
}

func NewHttpServer(host, port string, routeSetters ...RouteSetter) *Server {

	r := mux.NewRouter()

	for _, rs := range routeSetters {
		rs.Set(r)
	}

	s := &Server{
		Srv: &http.Server{
			Addr: fmt.Sprintf("%s:%s", host, port),
			Handler: r,
			WriteTimeout: time.Hour,
			ReadTimeout: time.Hour,
		},
	}

	return s
}

func (s *Server) Open() error {
	return s.Srv.ListenAndServe()
}

