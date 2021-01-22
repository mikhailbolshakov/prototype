package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"time"
)

type DeliveryService interface {
	// delivery user service
	Delivery(rq *DeliveryRequest) (*Delivery, error)
	Complete(deliveryId string, finishTime *time.Time) (*Delivery, error)
	Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error)
	Get(deliveryId string) *Delivery
	UpdateDetails(deliveryId string, details map[string]interface{}) (*Delivery, error)
}

func NewDeliveryService(balanceService UserBalanceService,
	userService users.Service,
	storage storage.Storage,
	queue queue.Queue,
	bpm bpm.Engine) DeliveryService {

	d := &deliveryServiceImpl{
		balance:     balanceService,
		userService: userService,
		storage:     storage,
		queue:       queue,
		bpm:         bpm,
	}

	return d
}

type deliveryServiceImpl struct {
	balance     UserBalanceService
	userService users.Service
	queue       queue.Queue
	storage     storage.Storage
	bpm         bpm.Engine
}

func (s *deliveryServiceImpl) userIdName(input string) string {

	if _, err := uuid.Parse(input); err == nil {
		return input
	} else {
		return s.userService.Get(input).Id
	}

}

func (d *deliveryServiceImpl) Delivery(rq *DeliveryRequest) (*Delivery, error) {

	userId := d.userIdName(rq.UserId)

	st, ok := d.balance.GetTypes()[rq.ServiceTypeId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	if userBalance, err := d.balance.Get(&GetBalanceRequest{UserId: userId}); err == nil {

		if b, ok := userBalance.Balance[st]; ok {

			if b.Available == 0 {
				return nil, errors.New(fmt.Sprintf("user %s has no available service of type %s; number of locked services is %d", rq.UserId, rq.ServiceTypeId, b.Locked))
			}

		} else {
			return nil, errors.New(fmt.Sprintf("user %s doesn't have service %s on balance", rq.UserId, rq.ServiceTypeId))
		}

	} else {
		return nil, err
	}

	// lock service
	if _, err := d.balance.Lock(&ModifyBalanceRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      1,
	}); err != nil {
		return nil, err
	}

	// create a delivery object
	delivery := &Delivery{
		Id:            kit.NewId(),
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Status:        "in-progress",
		StartTime:     time.Now().UTC(),
		FinishTime:    nil,
		Details:       rq.Details,
	}

	// save to storage
	dto, err := d.storage.CreateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}
	delivery = d.deliveryFromDto(dto)

	// execute a handler for corresponding task
	if st, ok := d.balance.GetTypes()[rq.ServiceTypeId]; ok {
		if st.DeliveryWfId != "" {
			// start WF
			variables := make(map[string]interface{})
			variables["deliveryId"] = delivery.Id
			_, err := d.bpm.StartProcess(st.DeliveryWfId, variables)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("WF process isn't specified for service type %s", rq.ServiceTypeId))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	// update storage
	dto, err = d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}
	delivery = d.deliveryFromDto(dto)

	return delivery, nil
}

func (d *deliveryServiceImpl) Get(deliveryId string) *Delivery{
	dto := d.storage.GetDelivery(deliveryId)
	return d.deliveryFromDto(dto)
}

func (d *deliveryServiceImpl) Complete(deliveryId string, finishTime *time.Time) (*Delivery, error) {

	delivery := d.deliveryFromDto(d.storage.GetDelivery(deliveryId))
	delivery.Status = "completed"
	if finishTime != nil {
		delivery.FinishTime = finishTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	dto, err := d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}

	_, err = d.balance.WriteOff(&ModifyBalanceRequest{
		UserId:        delivery.UserId,
		ServiceTypeId: delivery.ServiceTypeId,
		Quantity:      1,
	})
	if err != nil {
		return nil, err
	}

	return d.deliveryFromDto(dto), nil

}

func (d *deliveryServiceImpl) Cancel(deliveryId string, cancelTime *time.Time) (*Delivery, error) {

	delivery := d.deliveryFromDto(d.storage.GetDelivery(deliveryId))
	delivery.Status = "canceled"
	if cancelTime != nil {
		delivery.FinishTime = cancelTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	dto, err := d.storage.UpdateDelivery(d.deliveryToDto(delivery))
	if err != nil {
		return nil, err
	}

	return d.deliveryFromDto(dto), nil
}

func (d *deliveryServiceImpl) UpdateDetails(deliveryId string, details map[string]interface{}) (*Delivery, error) {

	dt := "{}"
	if details != nil {
		bytes, _ := json.Marshal(details)
		dt = string(bytes)
	}

	dto, err := d.storage.UpdateDetails(deliveryId, dt)
	if err != nil {
		return nil, err
	}

	return d.deliveryFromDto(dto), nil
}

