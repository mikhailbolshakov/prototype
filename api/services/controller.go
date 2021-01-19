package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
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

func (c *controller) AddBalance(writer http.ResponseWriter, request *http.Request) {

	rq := &ModifyUserBalanceRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userId := mux.Vars(request)["userId"]

	if rsPb, err := c.grpc.balance.Add(ctx, &pb.ChangeServicesRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int32(rq.Quantity),
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.balanceFromPb(rsPb))
	}

}

func (c *controller) GetBalance(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userId := mux.Vars(request)["userId"]

	if rsPb, err := c.grpc.balance.GetBalance(ctx, &pb.GetBalanceRequest{
		UserId:        userId,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.balanceFromPb(rsPb))
	}
}

func (c *controller) Delivery(writer http.ResponseWriter, request *http.Request) {

	rq := &DeliveryRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userId := mux.Vars(request)["userId"]

	detailsB, err := json.Marshal(rq.Details)
	if err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	if rsPb, err := c.grpc.delivery.Create(ctx, &pb.DeliveryRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Details:       detailsB,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.deliveryFromPb(rsPb))
	}

}