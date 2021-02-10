package services

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	pb "gitlab.medzdrav.ru/prototype/proto/services"
	"net/http"
)

type Controller interface {
	AddBalance(http.ResponseWriter, *http.Request)
	GetBalance(http.ResponseWriter, *http.Request)
	Delivery(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	balanceService public.BalanceService
	deliveryService public.DeliveryService
}

func NewController(balanceService public.BalanceService, deliveryService public.DeliveryService) Controller {
	return &ctrlImpl{
		balanceService: balanceService,
		deliveryService: deliveryService,
	}
}

func (c *ctrlImpl) AddBalance(writer http.ResponseWriter, request *http.Request) {

	rq := &ModifyUserBalanceRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	userId := mux.Vars(request)["userId"]

	if rsPb, err := c.balanceService.Add(&pb.ChangeServicesRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int32(rq.Quantity),
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.balanceFromPb(rsPb))
	}

}

func (c *ctrlImpl) GetBalance(writer http.ResponseWriter, request *http.Request) {

	userId := mux.Vars(request)["userId"]

	if rsPb, err := c.balanceService.GetBalance(&pb.GetBalanceRequest{
		UserId:        userId,
	}); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.balanceFromPb(rsPb))
	}
}

func (c *ctrlImpl) Delivery(writer http.ResponseWriter, request *http.Request) {

	rq := &DeliveryRequest{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	userId := mux.Vars(request)["userId"]

	if rsPb, err := c.deliveryService.Create(userId, rq.ServiceTypeId, rq.Details); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, c.deliveryFromPb(rsPb))
	}

}