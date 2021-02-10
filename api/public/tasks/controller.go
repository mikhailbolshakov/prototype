package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/tasks"
	"net/http"
	"strconv"
	"time"
)

type Controller interface {
	New(http.ResponseWriter, *http.Request)
	MakeTransition(http.ResponseWriter, *http.Request)
	GetById(http.ResponseWriter, *http.Request)
	SetAssignee(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	AssignmentLog(http.ResponseWriter, *http.Request)
	GetHistory(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	taskService public.TaskService
}

func NewController(taskService public.TaskService) Controller {
	return &ctrlImpl{
		taskService: taskService,
	}
}

func (c *ctrlImpl) New(writer http.ResponseWriter, request *http.Request) {

	rq := &NewTaskRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.taskService.New(c.toPb(rq)); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) MakeTransition(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]
	transitionId := mux.Vars(request)["transitionId"]

	if t, err := c.taskService.MakeTransition(&pb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: transitionId,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(t))
	}
}

func (c *ctrlImpl) GetById(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]

	if rsPb, err := c.taskService.GetById(taskId); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}
}

func (c *ctrlImpl) SetAssignee(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]

	rq := &Assignee{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.taskService.SetAssignee(&pb.SetAssigneeRequest{
		TaskId: taskId,
		Assignee: &pb.Assignee{
			Type:     rq.Type,
			Group:    rq.Group,
			UserId:   rq.UserId,
			Username: rq.Username,
			At:       grpc.TimeToPbTS(rq.At),
		},
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}
}

func (c *ctrlImpl) Search(writer http.ResponseWriter, request *http.Request) {

	rq := &pb.SearchRequest{
		Paging: &pb.PagingRequest{},
		Status: &pb.Status{
			Status:    request.FormValue("status"),
			Substatus: request.FormValue("substatus"),
		},
		Assignee: &pb.Assignee{
			Group:    request.FormValue("group"),
			Username: request.FormValue("username"),
			UserId:   request.FormValue("userId"),
		},
		Type: &pb.Type{
			Type:    request.FormValue("type"),
			Subtype: request.FormValue("subtype"),
		},
		Num: request.FormValue("num"),
	}

	if sizeTxt := request.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("limit: "+e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := request.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("offset: "+e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	if rsPb, err := c.taskService.Search(rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.searchRsFromPb(rsPb))
	}

}

func (c *ctrlImpl) AssignmentLog(writer http.ResponseWriter, request *http.Request) {

	rq := &pb.AssignmentLogRequest{
		Paging: &pb.PagingRequest{},
	}

	if startBeforeTxt := request.FormValue("startBefore"); startBeforeTxt != "" {
		startBefore, e := time.Parse(time.RFC3339, startBeforeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("startBefore: "+e.Error()))
			return
		}
		rq.StartTimeBefore = grpc.TimeToPbTS(&startBefore)
	}

	if startAfterTxt := request.FormValue("startAfter"); startAfterTxt != "" {
		startAfter, e := time.Parse(time.RFC3339, startAfterTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("startAfter: "+e.Error()))
			return
		}
		rq.StartTimeAfter = grpc.TimeToPbTS(&startAfter)
	}

	if sizeTxt := request.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("limit: "+e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := request.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("offset: "+e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	if rsPb, err := c.taskService.GetAssignmentLog(rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.assLogRsFromPb(rsPb))
	}

}

func (c *ctrlImpl) GetHistory(writer http.ResponseWriter, request *http.Request) {

	taskId := mux.Vars(request)["id"]

	if rsPb, err := c.taskService.GetHistory(taskId); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.histFromPb(rsPb))
	}
}
