package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"net/http"
	"strconv"
	"time"
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

	if rsPb, err := c.grpc.tasks.New(ctx, c.toPb(rq)); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
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
		c.RespondOK(writer, c.fromPb(rsPb))
	}
}

func (c *controller) GetById(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.tasks.GetById(ctx, &pb.GetByIdRequest{
		Id:       taskId,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}
}

func (c *controller) SetAssignee(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]

	rq := &Assignee{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.tasks.SetAssignee(ctx, &pb.SetAssigneeRequest{
		TaskId:   taskId,
		Assignee: &pb.Assignee{
			Group: rq.Group,
			User:  rq.User,
			At:    grpc.TimeToPbTS(rq.At),
		},
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}
}

func (c *controller) Search(writer http.ResponseWriter, request *http.Request) {

	rq := &pb.SearchRequest{
		Paging:   &pb.PagingRequest{},
		Status:   &pb.Status{
			Status:    request.FormValue("status"),
			Substatus: request.FormValue("substatus"),
		},
		Assignee: &pb.Assignee{
			Group: request.FormValue("group"),
			User:  request.FormValue("user"),
		},
		Type:     &pb.Type{
			Type:    request.FormValue("type"),
			Subtype: request.FormValue("subtype"),
		},
	}

	if sizeTxt := request.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("limit: " + e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := request.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("offset: " + e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.tasks.Search(ctx, rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.searchRsFromPb(rsPb))
	}

}

func (c *controller) AssignmentLog(writer http.ResponseWriter, request *http.Request) {

	rq := &pb.AssignmentLogRequest{
		Paging:   &pb.PagingRequest{},
	}

	if startBeforeTxt := request.FormValue("startBefore"); startBeforeTxt != "" {
		startBefore, e := time.Parse(time.RFC3339, startBeforeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("startBefore: " + e.Error()))
			return
		}
		rq.StartTimeBefore = grpc.TimeToPbTS(&startBefore)
	}

	if startAfterTxt := request.FormValue("startAfter"); startAfterTxt != "" {
		startAfter, e := time.Parse(time.RFC3339, startAfterTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("startAfter: " + e.Error()))
			return
		}
		rq.StartTimeAfter = grpc.TimeToPbTS(&startAfter)
	}

	if sizeTxt := request.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("limit: " + e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := request.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("offset: " + e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.tasks.GetAssignmentLog(ctx, rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.assLogRsFromPb(rsPb))
	}

}