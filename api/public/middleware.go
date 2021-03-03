package public

import (
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"net/http"
)

const (
	HTTP_HEADER_SESSION_ID = "X-SESSION-ID"
)

type mdw struct {
	sessionsService SessionsService
}

func NewMiddleware(sessionsService SessionsService) *mdw {
	return &mdw{sessionsService: sessionsService}
}

func (m *mdw) SessionMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		sessionId := r.Header.Get(HTTP_HEADER_SESSION_ID)

		if sessionId == "" {
			http.Error(w, "session id missing", http.StatusUnauthorized)
			return
		}

		ctxRq := kitContext.NewRequestCtx().
			Rest().
			WithNewRequestId().
			WithSessionId(sessionId)

		ctx := ctxRq.ToContext(r.Context())

		session, err := m.sessionsService.AuthSession(ctx, sessionId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if session == nil {
			http.Error(w, "no active session", http.StatusUnauthorized)
			return
		}

		ctx = ctxRq.
			WithUser(session.UserId, session.Username).
			WithChatUserId(session.ChatUserId).
			ToContext(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (m *mdw) NoSessionMiddleware(next http.Handler) http.Handler {

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