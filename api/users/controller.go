package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/kit/http/auth"
	pb "gitlab.medzdrav.ru/prototype/proto/users"
	"net/http"
	"strconv"
	"strings"
)

type controller struct {
	kitHttp.Controller
	grpc *grpcClient
	auth auth.AuthenticationHandler
}

func newController(auth auth.AuthenticationHandler) (*controller, error) {

	c, err := newGrpcClient()
	if err != nil {
		return nil, err
	}

	return &controller{
		auth: auth,
		grpc: c,
	}, nil
}

func (c *controller) Create(writer http.ResponseWriter, request *http.Request) {

	rq := &CreateUserRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.users.Create(ctx, c.toPb(rq)); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *controller) GetByUsername(writer http.ResponseWriter, request *http.Request) {

	username := mux.Vars(request)["un"]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.users.GetByUsername(ctx, &pb.GetByUsernameRequest{Username: username}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {

		if rsPb == nil {
			c.RespondError(writer, http.StatusNotFound, errors.New("user not found"))
		}

		c.RespondOK(writer, c.fromPb(rsPb))
	}

}

func (c *controller) Search(writer http.ResponseWriter, request *http.Request) {

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.users.Search(ctx, rq); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.searchRsFromPb(rsPb))
	}

}

func (c *controller) Login(writer http.ResponseWriter, request *http.Request) {

	rq := &LoginRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	jwt, err := c.auth.AuthenticateUser(&auth.AuthenticateUser{
		UserName: rq.Username,
		Password: rq.Password,
	})
	if err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
		return
	}

	rs := &LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

	c.RespondOK(writer, rs)

}
