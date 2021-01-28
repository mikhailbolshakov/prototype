package session

import (
	"fmt"
	"net/http"
)

const (
	HTTP_HEADER_SESSION_ID = "X-SESSION-ID"
)

func (h *hubImpl) SessionMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		sessionId := r.Header.Get(HTTP_HEADER_SESSION_ID)

		if sessionId == "" {
			http.Error(w, "session id missing", http.StatusUnauthorized)
			return
		}

		session := h.GetById(sessionId)

		if session == nil {
			http.Error(w, "no active session", http.StatusUnauthorized)
			return
		}

		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.getAccessToken()))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}