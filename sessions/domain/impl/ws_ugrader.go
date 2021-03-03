package impl

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"net/http"
)

type upgrader struct {
	hub Hub
}

func newWsUpgrader(hub Hub) http2.WsUpgrader {
	return &upgrader{
		hub: hub,
	}
}

func (u *upgrader) Set(noAuthRouter *mux.Router, upgrader *websocket.Upgrader) {

	noAuthRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		l := log.L().Pr("ws").Cmp("mdw").Mth("upgrade")

		w.Header().Set("Content-Type", "application/json")
		wsConn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			l.E(err).Err()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionId := r.URL.Query().Get("session")
		l.F(log.FF{"sid": sessionId}).Inf("request")

		if sessionId == "" {
			l.Err("no session provided")
			http.Error(w, "no session provided", http.StatusBadRequest)
			return
		}

		if err := u.hub.setupWsConnection(sessionId, wsConn); err != nil {
			l.E(err).Err()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		l.Inf("ok ")

	})

}
