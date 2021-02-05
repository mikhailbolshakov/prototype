package session

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"time"
)

const (
	writeWait           = 1000 * time.Second
	pongWait            = 6000 * time.Second
	pingPeriod          = (pongWait * 9) / 10
	maxMessageSize      = 4608
)

type Ws interface {
	open()
	close()
	closedEvent() <-chan struct{}
	receivedMessageEvent() <-chan []byte
	send(message []byte)
}

type wsImpl struct {
	sessionId    string
	userId       string
	conn         *websocket.Conn
	sendChan     chan []byte
	receivedChan chan []byte
	closeChan    chan struct{}
}

func newWs(conn *websocket.Conn, sessionId, userId string) Ws {
	return &wsImpl{
		conn:         conn,
		sendChan:     make(chan []byte, 256),
		receivedChan: make(chan []byte, 256),
		closeChan:    make(chan struct{}),
		sessionId:    sessionId,
		userId:       userId,
	}
}

func (w *wsImpl) closedEvent() <-chan struct{} {
	return w.closeChan
}

func (w *wsImpl) receivedMessageEvent() <-chan []byte {
	return w.receivedChan
}

func (w *wsImpl) open() {
	go w.write()
	go w.read()
}

func (w *wsImpl) close() {
	w.closeChan <- struct{}{}
	_ = w.conn.Close()
}

func (w *wsImpl) send(message []byte) {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				log.Err(err, true)
			}
		}
	}()

	w.sendChan <- message
}

func (w *wsImpl) write() {

	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		// close connection
		close(w.sendChan)
		pingTicker.Stop()
		w.close()
	}()

	for {
		select {

		case message, ok := <-w.sendChan:

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Dbg("[ws] close message sent")
				return
			}

			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Err(err, true)
				return
			}

			_, err = writer.Write(message)
			if err != nil {
				log.Err(err, true)
				return
			}
			log.Dbg("[ws] message sent: ", string(message))

			if err := writer.Close(); err != nil {
				log.Err(err, true)
				return
			}

		case <-pingTicker.C:
			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Err(err, true)
				return
			}
		}
	}
}

func (w *wsImpl) read() {

	defer func() {
		w.close()
	}()

	w.conn.SetReadLimit(maxMessageSize)
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := w.conn.ReadMessage()
		log.TrcF("[ws] message received %s", string(message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Err(err, true)
			} else {
				log.Err(fmt.Errorf("read socket error: %s", err.Error()), true)
			}
			return
		}

		if string(message) == "ping" {
			w.sendChan <- []byte("pong")
		} else {
			w.receivedChan <- message
		}

	}
}
