package api

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func (h *TestHelper) WebrtcWs(sessionId, roomId string) (*websocket.Conn, chan struct{}, error) {

	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?session=%s&room=%s", WEBRTC_URL, sessionId, roomId), nil)
	if err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	go h.ListenWs(wsConn, done)
	return wsConn, done, nil
}
