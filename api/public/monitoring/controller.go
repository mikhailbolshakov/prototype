package monitoring

import (
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Controller interface {
	GetUserSessions(http.ResponseWriter, *http.Request)
	GetTotalSessions(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	monitor public.SessionMonitor
}

func NewController(monitor public.SessionMonitor) Controller {
	return &ctrlImpl{
		monitor: monitor,
	}
}

func (c *ctrlImpl) GetUserSessions(writer http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	rs, err := c.monitor.UserSessions(r.Context(), userId)
	if err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
		return
	}
	c.RespondOK(writer, c.toUserSessionsApi(rs))
}

func (c *ctrlImpl) GetTotalSessions(writer http.ResponseWriter, r *http.Request) {
	rs, err := c.monitor.TotalSessions(r.Context())
	if err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
		return
	}
	c.RespondOK(writer, c.toTotalSessionsApi(rs))
}
