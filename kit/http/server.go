package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"net/http"
	"time"
)

type Server struct {
	Srv          *http.Server
	AuthRouter   *mux.Router
	NoAuthRouter *mux.Router
	WsUpgrader   *websocket.Upgrader
	Cert, Key    string
	logger       log.CLoggerFunc
}

type RouteSetter interface {
	Set(authRouter, noAuthRouter *mux.Router)
}

type WsUpgrader interface {
	Set(noAuthRouter *mux.Router, upgrader *websocket.Upgrader)
}

func NewHttpServer(host, port, cert, key string, logger log.CLoggerFunc) *Server {

	r := mux.NewRouter()
	noAuthRouter := r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.Header.Get("Authorization") == ""
	}).Subrouter()
	authRouter := r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}).Subrouter()

	s := &Server{
		Srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      r,
			WriteTimeout: time.Hour,
			ReadTimeout:  time.Hour,
		},
		WsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		AuthRouter:   authRouter,
		NoAuthRouter: noAuthRouter,
		Cert:         cert,
		Key:          key,
		logger:       logger,
	}

	return s
}

func (s *Server) SetRouters(routeSetters ...RouteSetter) {
	for _, rs := range routeSetters {
		rs.Set(s.AuthRouter, s.NoAuthRouter)
	}
}

func (s *Server) SetWsUpgrader(upgradeSetter WsUpgrader) {
	upgradeSetter.Set(s.NoAuthRouter, s.WsUpgrader)
}

func (s *Server) SetAuthMiddleware(mdws ...mux.MiddlewareFunc) {
	for _, m := range mdws {
		s.AuthRouter.Use(m)
	}
}

func (s *Server) SetNoAuthMiddleware(mdws ...mux.MiddlewareFunc) {
	for _, m := range mdws {
		s.NoAuthRouter.Use(m)
	}
}

func (s *Server) SetMiddleware(mdws ...mux.MiddlewareFunc) {
	for _, m := range mdws {
		s.NoAuthRouter.Use(m)
		s.AuthRouter.Use(m)
	}
}

func (s *Server) Listen() {
	go func() {

		l := s.logger().Pr("http").Cmp("server").Mth("listen").F(log.FF{"url": s.Srv.Addr})
		l.Inf("start listening")

		// if tls parameters are specified, list tls connection
		if s.Cert == "" || s.Key == "" {
			l.E(s.Srv.ListenAndServe()).Err()
		} else {
			l.E(s.Srv.ListenAndServeTLS(s.Cert, s.Key)).Err()
		}
	}()
}

func (s *Server) Close() {
	_ = s.Srv.Close()
}
