package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"net/http"
)

type controller struct {
	kitHttp.Controller
	grpc *grpcClient
}

func newController() (*controller, error) {

	c, err := newGrpcClient()
	if err != nil {
		return nil, err
	}

	return &controller{
		grpc: c,
	}, nil
}

func (c *controller) New(writer http.ResponseWriter, request *http.Request) {

	rq := &NewTaskRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pb, err := c.toPb(rq)
	if err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.grpc.tasks.New(ctx, pb); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		rs, err := c.fromPb(rsPb)
		if err != nil {
			c.RespondError(writer, http.StatusInternalServerError, err)
		}
		c.RespondOK(writer, rs)
	}

}

func (c *controller) MakeTransition(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]
	transitionId := mux.Vars(request)["transitionId"]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.tasks.MakeTransition(ctx, &pb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: transitionId,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		rs, err := c.fromPb(rsPb)
		if err != nil {
			c.RespondError(writer, http.StatusInternalServerError, err)
		}
		c.RespondOK(writer, rs)
	}
}