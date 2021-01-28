package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"net/http"
	"strconv"
	"strings"
)

type Controller interface {
	CreateClient(http.ResponseWriter, *http.Request)
	CreateConsultant(http.ResponseWriter, *http.Request)
	CreateExpert(http.ResponseWriter, *http.Request)
	GetByUsername(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	userService users.Service
	auth        auth.Service
}

func NewController(auth auth.Service, userService users.Service) Controller {

	return &ctrlImpl{
		auth: auth,
		userService: userService,
	}
}

func (c *ctrlImpl) CreateClient(writer http.ResponseWriter, request *http.Request) {

	rq := &CreateClientRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	p := &pb.CreateClientRequest{
		FirstName:  rq.FirstName,
		MiddleName: rq.MiddleName,
		LastName:   rq.LastName,
		Sex:        rq.Sex,
		BirthDate:  grpc.TimeToPbTS(&rq.BirthDate),
		Phone:      rq.Phone,
		Email:      rq.Email,
	}

	if rsPb, err := c.userService.CreateClient(p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) CreateConsultant(writer http.ResponseWriter, request *http.Request) {

	rq := &CreateConsultantRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	p := &pb.CreateConsultantRequest{
		FirstName:  rq.FirstName,
		MiddleName: rq.MiddleName,
		LastName:   rq.LastName,
		Email:      rq.Email,
	}

	if rsPb, err := c.userService.CreateConsultant(p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) CreateExpert(writer http.ResponseWriter, request *http.Request) {

	rq := &CreateExpertRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	p := &pb.CreateExpertRequest{
		FirstName:      rq.FirstName,
		MiddleName:     rq.MiddleName,
		LastName:       rq.LastName,
		Email:          rq.Email,
		Specialization: rq.Specialization,
	}

	if rsPb, err := c.userService.CreateExpert(p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) GetByUsername(writer http.ResponseWriter, request *http.Request) {

	username := mux.Vars(request)["un"]

	if usr := c.userService.Get(username); usr != nil {
		c.RespondOK(writer, c.fromPb(usr))
	} else {
		c.RespondError(writer, http.StatusNotFound, errors.New("user not found"))
	}

}

func (c *ctrlImpl) Search(writer http.ResponseWriter, request *http.Request) {

	rq := &pb.SearchRequest{
		Paging: &pb.PagingRequest{
			Size:  0,
			Index: 0,
		},
		UserType:       request.FormValue("type"),
		Username:       request.FormValue("username"),
		Email:          request.FormValue("email"),
		Phone:          request.FormValue("phone"),
		MMId:           request.FormValue("mmId"),
		MMChannelId:    request.FormValue("channel"),
		OnlineStatuses: []string{},
	}

	if onlineStatusesTxt := request.FormValue("statuses"); onlineStatusesTxt != "" {
		rq.OnlineStatuses = strings.Split(onlineStatusesTxt, ",")
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

	if rsPb, err := c.userService.Search(rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.searchRsFromPb(rsPb))
	}

}

