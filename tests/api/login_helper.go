package api

import (
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/api/session"
)

func (h *TestHelper) Login(username string) (string, error) {

	rq := &session.LoginRequest{
		Username: username,
		Password: DEFAULT_PWD,
	}

	rqJ, _ := json.Marshal(rq)

	r, err := h.POST(fmt.Sprintf("%s/api/users/login", BASE_URL), rqJ)
	if err != nil {
		return "", err
	} else {

		var rs *session.LoginResponse
		err = json.Unmarshal(r, &rs)
		if err != nil {
			return "", err
		}

		fmt.Printf("user %s logged in. session_Id=%s\n", username, rs.SessionId)

		h.sessionId = rs.SessionId

		return rs.SessionId, nil
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
