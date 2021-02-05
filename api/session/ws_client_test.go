package session

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"
)

func listen(c *websocket.Conn, done chan struct{}) {

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)


	//readMessageChan := make(chan []byte)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Println(string(message))
			//readMessageChan <- message
		}
	}()

	for {
		select {
		case <-done:
			c.Close()
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				log.Println("write:", err)
				return
			}
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
				//case <-time.After(time.Second):
			}
			return
		}
	}

}

func Test_ReconnectWithSameSession(t *testing.T) {

	done := make(chan struct{})

	rq := struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "79032220402",
		Password: "12345",
	}

	rqJ, _ := json.Marshal(rq)

	rs, err := http.Post("http://localhost:8000/api/users/login", "application/json", bytes.NewBuffer(rqJ))
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = json.NewDecoder(rs.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(data)

	sessionId := data["sessionId"].(string)

	header := http.Header{}
	c, _, err := websocket.DefaultDialer.Dial( "ws://localhost:8000/ws?session=" + sessionId, header)
	if err != nil {
		t.Fatal(err)
	}

	go listen(c, done)

	<-time.After(time.Second * 20)
	done <- struct{}{}

	time.Sleep(time.Second)

	c, _, err = websocket.DefaultDialer.Dial( "ws://localhost:8000/ws?session=" + sessionId, header)
	if err != nil {
		t.Fatal(err)
	}

	listen(c, done)

}

func Test_MultipleConnectionsSameUser(t *testing.T) {

	done := make(chan struct{})

	rq := struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "79032220402",
		Password: "12345",
	}

	rqJ, _ := json.Marshal(rq)

	rs, err := http.Post("http://localhost:8000/api/users/login", "application/json", bytes.NewBuffer(rqJ))
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = json.NewDecoder(rs.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(data)

	sessionId1 := data["sessionId"].(string)

	header := http.Header{}
	c1, _, err := websocket.DefaultDialer.Dial( "ws://localhost:8000/ws?session=" + sessionId1, header)
	if err != nil {
		t.Fatal(err)
	}

	go listen(c1, done)


	rs, err = http.Post("http://localhost:8000/api/users/login", "application/json", bytes.NewBuffer(rqJ))
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(rs.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(data)

	sessionId2 := data["sessionId"].(string)

	header = http.Header{}
	c2, _, err := websocket.DefaultDialer.Dial( "ws://localhost:8000/ws?session=" + sessionId2, header)
	if err != nil {
		t.Fatal(err)
	}

	go listen(c2, done)

	<-time.After(time.Second * 30)
	done <- struct{}{}

}