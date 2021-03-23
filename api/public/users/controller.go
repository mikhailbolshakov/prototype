package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	"gitlab.medzdrav.ru/prototype/kit/grpc"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	sessionPb "gitlab.medzdrav.ru/prototype/proto/sessions"
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
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	sessionsService public.SessionsService
	userService     public.UserService
}

func NewController(sessionsService public.SessionsService, userService public.UserService) Controller {

	return &ctrlImpl{
		userService: userService,
		sessionsService: sessionsService,
	}
}

func (c *ctrlImpl) Login(w http.ResponseWriter, r *http.Request) {

	rq := &LoginRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	sid, err := c.sessionsService.Login(r.Context(), &sessionPb.LoginRequest{
		Username: rq.Username,
		Password: rq.Password,
		ChatLogin: rq.ChatLogin,
	})

	if err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	c.RespondOK(w, &LoginResponse{
		SessionId: sid,
	})
}

func (c *ctrlImpl) Logout(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["userId"]

	if userId == "" {
		c.RespondError(w, http.StatusBadRequest, errors.New("no user specified"))
		return
	}

	err := c.sessionsService.Logout(r.Context(), userId)
	if err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	c.RespondOK(w, &struct{}{})
}

func (c *ctrlImpl) CreateClient(writer http.ResponseWriter, r *http.Request) {

	rq := &CreateClientRequest{}
	decoder := json.NewDecoder(r.Body)
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
		PhotoUrl:   rq.PhotoUrl,
	}

	if rsPb, err := c.userService.CreateClient(r.Context(), p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) CreateConsultant(writer http.ResponseWriter, r *http.Request) {

	rq := &CreateConsultantRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	p := &pb.CreateConsultantRequest{
		FirstName:  rq.FirstName,
		MiddleName: rq.MiddleName,
		LastName:   rq.LastName,
		Email:      rq.Email,
		Groups:     rq.Groups,
		PhotoUrl:   rq.PhotoUrl,
	}

	if rsPb, err := c.userService.CreateConsultant(r.Context(), p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) CreateExpert(writer http.ResponseWriter, r *http.Request) {

	rq := &CreateExpertRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	p := &pb.CreateExpertRequest{
		FirstName:  rq.FirstName,
		MiddleName: rq.MiddleName,
		LastName:   rq.LastName,
		Email:      rq.Email,
		PhotoUrl:   rq.PhotoUrl,
		Groups:     rq.Groups,
	}

	if rsPb, err := c.userService.CreateExpert(r.Context(), p); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *ctrlImpl) GetByUsername(writer http.ResponseWriter, r *http.Request) {

	username := mux.Vars(r)["un"]

	if usr := c.userService.Get(r.Context(), username); usr != nil {
		c.RespondOK(writer, c.fromPb(usr))
	} else {
		c.RespondError(writer, http.StatusNotFound, errors.New("user not found"))
	}

}

func (c *ctrlImpl) Search(writer http.ResponseWriter, r *http.Request) {

	rq := &pb.SearchRequest{
		Paging: &pb.PagingRequest{
			Size:  0,
			Index: 0,
		},
		UserType:        r.FormValue("type"),
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Phone:           r.FormValue("phone"),
		MMId:            r.FormValue("mmId"),
		CommonChannelId: r.FormValue("commonChannel"),
		MedChannelId:    r.FormValue("medChannel"),
		LawChannelId:    r.FormValue("lawChannel"),
		UserGroup:       r.FormValue("group"),
		OnlineStatuses:  []string{},
	}

	if onlineStatusesTxt := r.FormValue("statuses"); onlineStatusesTxt != "" {
		rq.OnlineStatuses = strings.Split(onlineStatusesTxt, ",")
	}

	if sizeTxt := r.FormValue("limit"); sizeTxt != "" {
		size, e := strconv.Atoi(sizeTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("limit: "+e.Error()))
			return
		}
		rq.Paging.Size = int32(size)
	}

	if indexTxt := r.FormValue("offset"); indexTxt != "" {
		index, e := strconv.Atoi(indexTxt)
		if e != nil {
			c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("offset: "+e.Error()))
			return
		}
		rq.Paging.Index = int32(index)
	}

	if rsPb, err := c.userService.Search(r.Context(), rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.searchRsFromPb(rsPb))
	}

}
