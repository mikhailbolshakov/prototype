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

func (c *ctrlImpl) AddBalance(w http.ResponseWriter, r *http.Request) {

	rq := &ModifyUserBalanceRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	userId := mux.Vars(r)["userId"]

	if rsPb, err := c.balanceService.Add(r.Context(), &pb.ChangeServicesRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      int32(rq.Quantity),
	}); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.balanceFromPb(rsPb))
	}

}

func (c *ctrlImpl) GetBalance(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["userId"]

	if rsPb, err := c.balanceService.GetBalance(r.Context(), &pb.GetBalanceRequest{
		UserId:        userId,
	}); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.balanceFromPb(rsPb))
	}
}

func (c *ctrlImpl) Delivery(w http.ResponseWriter, r *http.Request) {

	rq := &DeliveryRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w,  http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	userId := mux.Vars(r)["userId"]

	if rsPb, err := c.deliveryService.Create(r.Context(), userId, rq.ServiceTypeId, rq.Details); err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(w, c.deliveryFromPb(rsPb))
	}

}