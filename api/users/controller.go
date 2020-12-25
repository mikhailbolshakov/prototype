package users

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"gitlab.medzdrav.ru/prototype/proto/users"
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

func (c *controller) Create(writer http.ResponseWriter, request *http.Request) {

	rq := &CreateUserRequest{}
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

	if rsPb, err := c.grpc.users.Create(ctx, pb); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		rs, err := c.fromPb(rsPb)
		if err != nil {
			c.RespondError(writer, http.StatusInternalServerError, err)
		}
		c.RespondOK(writer, rs)
	}

}

func (c *controller) GetByUsername(writer http.ResponseWriter, request *http.Request) {

	username := mux.Vars(request)["un"]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rsPb, err := c.grpc.users.GetByUsername(ctx, &users.GetByUsernameRequest{Username: username}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {

		if rsPb == nil {
			c.RespondError(writer, http.StatusNotFound, errors.New("user not found"))
		}

		rs, err := c.fromPb(rsPb)
		if err != nil {
			c.RespondError(writer, http.StatusInternalServerError, err)
		}
		c.RespondOK(writer, rs)
	}

}