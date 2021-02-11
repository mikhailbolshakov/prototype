package bp

import (
	"encoding/json"
	"errors"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/bp"
	"net/http"
)

type Controller interface {
	StartProcess(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	bpService public.BpService
}

func NewController(bpService public.BpService) Controller {
	return &ctrlImpl{
		bpService: bpService,
	}
}

func (c *ctrlImpl) StartProcess(w http.ResponseWriter, r *http.Request) {

	rq := &StartProcessRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	var varsB []byte
	if rq.Vars != nil {
		varsB, _ = json.Marshal(rq.Vars)
	}

	if rsPb, err := c.bpService.StartProcess(r.Context(), &pb.StartProcessRequest{
		ProcessId: rq.ProcessId,
		Vars:      varsB,
	}); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, &StartProcessResponse{Id: rsPb.Id})
	}

}
