package session

import (
	"fmt"
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

		w.Header().Set("Content-Type", "application/json")
		wsConn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Err(err, true)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionId := r.URL.Query().Get("session")
		if sessionId == "" {
			log.Err(fmt.Errorf("no session provided"), true)
			http.Error(w, "no session provided", http.StatusBadRequest)
			return
		}

		if err := u.hub.SetupWsConnection(sessionId, wsConn); err != nil {
			log.Err(err, true)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Dbg("[ws] session connected ", sessionId)

	})

}
