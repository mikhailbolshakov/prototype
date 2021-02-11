package session

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AuthCode string `json:"authCode"`
}

type LoginResponse struct {
	SessionId string `json:"sessionId"`
}

func (h *hubImpl) GetLoginRouteSetter() kitHttp.RouteSetter {
	return h
}

func (h *hubImpl) Set(authRouter, noAuthRouter *mux.Router) {

	noAuthRouter.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {

		rq := &LoginRequest{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(rq); err != nil {
			h.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
			return
		}

		rs, err := h.NewSession(r.Context(), &NewSessionRequest{
			Username:  rq.Username,
			Password:  rq.Password,
			ChatLogin: true,
		})
		if err != nil {
			h.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		h.RespondOK(w, &LoginResponse{
			SessionId: rs.SessionId,
		})

	}).Methods("POST")

	authRouter.HandleFunc("/api/users/{userId}/logout", func(w http.ResponseWriter, r *http.Request) {

		userId := mux.Vars(r)["userId"]

		if userId == "" {
			h.RespondError(w, http.StatusBadRequest, errors.New("no user specified"))
			return
		}

		if err := h.Logout(userId); err != nil {
			h.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		h.RespondOK(w, &struct{}{})

	}).Methods("POST")

}
