package monitoring

import (
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/session"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Controller interface {
	GetUserSessions(http.ResponseWriter, *http.Request)
	GetTotalSessions(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	monitor session.SessionMonitor
}

func NewController(monitor session.SessionMonitor) Controller {
	return &ctrlImpl{
		monitor: monitor,
	}
}

func (c *ctrlImpl) GetUserSessions(writer http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	rs := c.monitor.GetUserSessions(r.Context(), userId)
	c.RespondOK(writer, c.toUserSessionsApi(rs))
}

func (c *ctrlImpl) GetTotalSessions(writer http.ResponseWriter, r *http.Request) {
	rs := c.monitor.GetTotalSessions(r.Context())
	c.RespondOK(writer, c.toTotalSessionsApi(rs))
}
