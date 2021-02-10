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

func (c *ctrlImpl) StartProcess(writer http.ResponseWriter, request *http.Request) {

	rq := &StartProcessRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	var varsB []byte
	if rq.Vars != nil {
		varsB, _ = json.Marshal(rq.Vars)
	}

	if rsPb, err := c.bpService.StartProcess(&pb.StartProcessRequest{
		ProcessId: rq.ProcessId,
		Vars:      varsB,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, &StartProcessResponse{Id: rsPb.Id})
	}

}
