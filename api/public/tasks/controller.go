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

func (c *ctrlImpl) New(w http.ResponseWriter, r *http.Request) {

	rq := &NewTaskRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.taskService.New(r.Context(), c.toTaskRqPb(rq)); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toTaskApi(rsPb))
	}

}

func (c *ctrlImpl) MakeTransition(w http.ResponseWriter, r *http.Request) {

	taskId := mux.Vars(r)["id"]
	transitionId := mux.Vars(r)["transitionId"]

	if t, err := c.taskService.MakeTransition(r.Context(), &pb.MakeTransitionRequest{
		TaskId:       taskId,
		TransitionId: transitionId,
	}); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toTaskApi(t))
	}
}

func (c *ctrlImpl) GetById(w http.ResponseWriter, r *http.Request) {

	taskId := mux.Vars(r)["id"]

	if rsPb, err := c.taskService.GetById(r.Context(), taskId); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toTaskApi(rsPb))
	}
}

func (c *ctrlImpl) SetAssignee(w http.ResponseWriter, r *http.Request) {

	taskId := mux.Vars(r)["id"]

	rq := &Assignee{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.taskService.SetAssignee(r.Context(), &pb.SetAssigneeRequest{
		TaskId: taskId,
		Assignee: &pb.Assignee{
			Type:     rq.Type,
			Group:    rq.Group,
			UserId:   rq.UserId,
			Username: rq.Username,
			At:       grpc.TimeToPbTS(rq.At),
		},
	}); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toTaskApi(rsPb))
	}
}

func (c *ctrlImpl) Search(w http.ResponseWriter, r *http.Request) {

	rq := &pb.SearchRequest{
		Paging: &pb.PagingRequest{},
		Status: &pb.Status{
			Status:    r.FormValue("status"),
			Substatus: r.FormValue("substatus"),
		},
		Assignee: &pb.Assignee{
			Group:    r.FormValue("group"),
			Username: r.FormValue("username"),
			UserId:   r.FormValue("userId"),
		},
		Type: &pb.Type{
			Type:    r.FormValue("type"),
			Subtype: r.FormValue("subtype"),
		},
		Num: r.FormValue("num"),
	}

	if sizeTxt := r.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("limit: "+e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := r.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("offset: "+e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	if rsPb, err := c.taskService.Search(r.Context(), rq); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toSrchRsApi(rsPb))
	}

}

func (c *ctrlImpl) AssignmentLog(w http.ResponseWriter, r *http.Request) {

	rq := &pb.AssignmentLogRequest{
		Paging: &pb.PagingRequest{},
	}

	if startBeforeTxt := r.FormValue("startBefore"); startBeforeTxt != "" {
		startBefore, e := time.Parse(time.RFC3339, startBeforeTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("startBefore: "+e.Error()))
			return
		}
		rq.StartTimeBefore = grpc.TimeToPbTS(&startBefore)
	}

	if startAfterTxt := r.FormValue("startAfter"); startAfterTxt != "" {
		startAfter, e := time.Parse(time.RFC3339, startAfterTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("startAfter: "+e.Error()))
			return
		}
		rq.StartTimeAfter = grpc.TimeToPbTS(&startAfter)
	}

	if sizeTxt := r.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("limit: "+e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := r.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(w, http.StatusBadRequest, fmt.Errorf("offset: "+e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	if rsPb, err := c.taskService.GetAssignmentLog(r.Context(), rq); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toAssgnLogRsApi(rsPb))
	}

}

func (c *ctrlImpl) GetHistory(w http.ResponseWriter, r *http.Request) {

	taskId := mux.Vars(r)["id"]

	if rsPb, err := c.taskService.GetHistory(r.Context(), taskId); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.toHistApi(rsPb))
	}
}
