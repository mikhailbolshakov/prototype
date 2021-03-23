package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.medzdrav.ru/prototype/api/public/monitoring"
	"gitlab.medzdrav.ru/prototype/api/public/users"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func (h *TestHelper) Login(username string) (string, chan struct{}, error) {

	rq := &users.LoginRequest{
		Username: username,
		Password: DEFAULT_PWD,
		ChatLogin: true,
	}

	rqJ, _ := json.Marshal(rq)

	r, err := h.POST(fmt.Sprintf("%s/api/users/login", BASE_URL), rqJ)
	if err != nil {
		return "", nil, err
	} else {

		var rs *users.LoginResponse
		err = json.Unmarshal(r, &rs)
		if err != nil {
			return "", nil, err
		}

		fmt.Printf("user %s logged in. session_Id=%s\n", username, rs.SessionId)

		h.sessionId = rs.SessionId

		ws, done, err := h.Ws(h.sessionId)
		if err != nil {
			return "", nil, err
		}

		h.ws = ws
		return rs.SessionId, done, nil
	}
}

func (h *TestHelper) Logout(username string) error {

	user, err := h.GetUser(username)
	if err != nil {
		return err
	}

	_, err = h.POST(fmt.Sprintf("%s/api/users/%s/logout", BASE_URL, user.Id), []byte{})
	if err != nil {
		return err
	} else {

		fmt.Printf("user %s logged out\n", username)

		h.sessionId = ""

		return nil
	}
}

func (h *TestHelper) Ws(sessionId string) (*websocket.Conn, chan struct{}, error) {

	header := http.Header{}
	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?session=%s", WS_URL, sessionId), header)
	if err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	go h.ListenWs(wsConn, done)
	return wsConn, done, nil
}

func (h *TestHelper) ListenWs(c *websocket.Conn, done chan struct{}) {

	ticker := time.NewTicker(h.wsPingInterval)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c.SetPongHandler(func(m string) error {
		log.Println("[ws] received pong")
		return nil
	})

	//readMessageChan := make(chan []byte)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("[ws] read error:", err)
				return
			}
			log.Println("[ws] received:", string(message))
			//readMessageChan <- message
		}
	}()

	for {
		select {
		case <-done:
			c.Close()
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Println("[ws] write error:", err)
				return
			}
			log.Println("[ws] send ping")
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			}
			return
		}
	}

}

func (h *TestHelper) MonitorUserSessions(userId string) (*monitoring.UserSessionInfo, error) {

	rs, err := h.GET(fmt.Sprintf("%s/api/monitor/sessions/users/%s", BASE_URL, userId))
	if err != nil {
		return nil, err
	}
	var res *monitoring.UserSessionInfo
	_ = json.Unmarshal(rs, &res)
	return res, nil

}
