package api

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BASE_URL    = "http://localhost:8000"
	WS_URL      = "ws://localhost:8000/ws"
	APPJSON     = "application/json"
	MM_URL      = "http://localhost:8065"
	DEFAULT_PWD = "12345"
	TEST_USER   = "79032221002"
)

type TestHelper struct {
	sessionId string
	ws *websocket.Conn
	wsPingInterval time.Duration
}

func NewTestHelper() *TestHelper {
	return &TestHelper{
		wsPingInterval: time.Second,
	}
}

func (h *TestHelper) do(url string, verb string, payload []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(verb, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if h.sessionId != "" {
		req.Header.Add("X-SESSION-ID", h.sessionId)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("response error: %s", string(data))
	}

	return data, nil

}

func (h *TestHelper) POST(url string, payload []byte) ([]byte, error) {
	return h.do(url, "POST", payload)
}

func (h *TestHelper) PUT(url string, payload []byte) ([]byte, error) {
	return h.do(url, "PUT", payload)
}

func (h *TestHelper) GET(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if h.sessionId != "" {
		req.Header.Add("X-SESSION-ID", h.sessionId)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("response error: %s", string(data))
	}

	return data, nil

}
