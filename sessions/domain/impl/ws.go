package impl

import (
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/sessions/logger"
	"go.uber.org/atomic"
	"net"
	"time"
)

const (
	writeWait      = 1000 * time.Second
	pongWait       = 6000 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4608
	keepAlive      = time.Minute
)

// if we got no ping messages for "keepAlive" period, we close WS connection

type Ws interface {
	open()
	close()
	wsClosedEvent() <-chan struct{}
	receivedMessageEvent() <-chan []byte
	send(message []byte)
}

type wsImpl struct {
	sessionId    string
	userId       string
	conn         *websocket.Conn
	sendChan     chan []byte
	receivedChan chan []byte
	wsClosedChan chan struct{}
	pingChan     chan struct{}
	closed       *atomic.Bool
}

func newWs(conn *websocket.Conn, sessionId, userId string) Ws {

	return &wsImpl{
		conn:         conn,
		sendChan:     make(chan []byte, 256),
		receivedChan: make(chan []byte, 256),
		wsClosedChan: make(chan struct{}),
		sessionId:    sessionId,
		userId:       userId,
		pingChan:     make(chan struct{}),
		closed:       atomic.NewBool(false),
	}
}

func (w *wsImpl) l() log.CLogger {
	return logger.L().Cmp("ws")
}

func (w *wsImpl) wsClosedEvent() <-chan struct{} {
	return w.wsClosedChan
}

func (w *wsImpl) receivedMessageEvent() <-chan []byte {
	return w.receivedChan
}

func (w *wsImpl) open() {
	go w.write()
	go w.read()
	go w.pingListener()
}

func (w *wsImpl) close() {
	if !w.closed.Load() {
		w.wsClosedChan <- struct{}{}
		_ = w.conn.Close()
		w.closed.Store(true)
	}
}

func (w *wsImpl) send(message []byte) {

	l := w.l().Mth("ws-send")
	l.TrcF("%s", string(message))

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				l.E(err).St().Err("recovered")
			}
		}
	}()

	w.sendChan <- message
}

func (w *wsImpl) write() {

	defer func() {
		// close connection
		close(w.sendChan)
		w.close()
	}()

	l := w.l().Mth("write")

	for {
		select {

		case message, ok := <-w.sendChan:

			_ = w.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				l.Dbg("close message sent")
				return
			}

			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				l.E(err).St().Err()
				return
			}

			_, err = writer.Write(message)
			if err != nil {
				l.E(err).St().Err()
				return
			}
			l.TrcF("message sent: %s", string(message))

			if err := writer.Close(); err != nil {
				l.E(err).St().Err()
				return
			}

		}
	}
}

func (w *wsImpl) pingListener() {
	l := w.l().Mth("ping-listener")
	for {
		select {
		case <-w.pingChan:
			if w.closed.Load() {
				return
			}
		case <-time.After(keepAlive):
			l.InfF("keep alive period elapsed, ws is closed")
			w.close()
			return
		}
	}
}

func (w *wsImpl) pingHandler(message string) error {
	w.pingChan <- struct{}{}
	err := w.conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(writeWait))
	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}

func (w *wsImpl) read() {

	l := w.l().Mth("read")

	defer func() {
		close(w.receivedChan)
		w.close()
	}()

	w.conn.SetReadLimit(maxMessageSize)
	_ = w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPingHandler(w.pingHandler)

	for {
		_, message, err := w.conn.ReadMessage()
		strMsg := string(message)
		l.TrcF("message received '%s'\n", strMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				l.E(err).Warn("socket closed")
			} else {
				l.E(err).Warn("read socket error")
			}
			w.wsClosedChan <- struct{}{}
			return
		}

		if strMsg == "" {
			l.Warn("empty message")
			continue
		}

		w.receivedChan <- message

	}

}
