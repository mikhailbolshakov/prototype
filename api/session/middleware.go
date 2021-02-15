package session

import (
	"fmt"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
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

		newCtx := kitContext.NewRequestCtx().
			Rest().
			WithNewRequestId().
			WithSessionId(sessionId).
			WithUser(session.getUserId(), session.getUsername()).
			WithChatUserId(session.getChatUserId()).
			ToContext(r.Context())

		r = r.WithContext(newCtx)

		// Currently we check token only on API level
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.getAccessToken()))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (h *hubImpl) NoSessionMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		newCtx := kitContext.NewRequestCtx().
			Rest().
			WithNewRequestId().
			ToContext(r.Context())

		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}