package session

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	http2 "gitlab.medzdrav.ru/prototype/kit/http"
	"log"
	"net/http"
	"time"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionId := r.URL.Query().Get("session")
		if sessionId == "" {
			http.Error(w, "no session provided", http.StatusBadRequest)
			return
		}

		if err := u.hub.SetWs(sessionId, newWs(wsConn)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4608
)

type Ws interface {
	open()
	close()
	send(message []byte)
}

type wsImpl struct {
	conn     *websocket.Conn
	sendChan chan []byte
}

func newWs(conn *websocket.Conn) Ws {
	return &wsImpl{
		conn: conn,
		sendChan: make(chan []byte, 256),
	}
}

func (w *wsImpl) open() {
	go w.write()
	go w.read()
}

func (w *wsImpl) close() {
	close(w.sendChan)
	_ = w.conn.Close()
}

func (w *wsImpl) send(message []byte) {

	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(error)
			if !ok {
				log.Println(fmt.Errorf("%v", r))
			}
		}
	}()

	w.sendChan <- message
}

func (w *wsImpl) write() {

	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		// close connection
		pingTicker.Stop()
		_ = w.conn.Close()
	}()

	for {
		select {

		case message, ok := <- w.sendChan:

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				//log.Printf("Channel has been closed for account %s", c.account.Id)
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			log.Println("message to client: ", string(message))

			_, _ = writer.Write(message)
			if err := writer.Close(); err != nil {
				log.Println("writer.Close() error:", err.Error())
				return
			}

		case <-pingTicker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("write ping message error:", err.Error())
				return
			}
		}
	}
}

func (w *wsImpl) read() {

	defer func() {
		//app.L().Debugf("Websocket client is closing (accountId: %s)", c.account.Id.String())
		//w.hub.unregisterChan <- c
		//w.conn.Close()
	}()

	w.conn.SetReadLimit(maxMessageSize)
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := w.conn.ReadMessage()
		//app.L().Debugf("Message from socket: %s", string(message))

		if err != nil {
			log.Println("read socket error:", err.Error())
			//app.E().SetError(system.SysErr(err, system.WsConnReadMessageErrorCode, message))
			//app.L().Debug(">>>>>> ReadMessageError:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//app.L().Debugf("> > > Read sentry: %s", err)
			}
			break
		}

		log.Println(message)

		//go w.hub.onMessage(message, c)
	}
}
