package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BASE_URL    = "http://localhost:8000"
	APPJSON     = "application/json"
	MM_URL      = "http://localhost:8065"
	DEFAULT_PWD = "12345"
	TEST_USER  = "79032221002"
)

type TestHelper struct {
	sessionId string
}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) POST (url string, payload []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
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