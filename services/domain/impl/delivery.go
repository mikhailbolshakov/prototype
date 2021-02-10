package impl

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"time"
)

func NewDeliveryService(balanceService domain.UserBalanceService,
	userService domain.UserService,
	bpService domain.BpService,
	storage domain.Storage,
	queue queue.Queue) domain.DeliveryService {

	d := &deliveryServiceImpl{
		balance:     balanceService,
		userService: userService,
		storage:     storage,
		queue:       queue,
		bpService:   bpService,
	}

	return d
}

type deliveryServiceImpl struct {
	balance     domain.UserBalanceService
	userService domain.UserService
	queue       queue.Queue
	storage     domain.Storage
	bpService   domain.BpService
}

func (s *deliveryServiceImpl) userIdName(input string) string {
	if _, err := uuid.Parse(input); err == nil {
		return input
	} else {
		return s.userService.Get(input).Id
	}
}

func (d *deliveryServiceImpl) Delivery(rq *domain.DeliveryRequest) (*domain.Delivery, error) {

	userId := d.userIdName(rq.UserId)

	st, ok := d.balance.GetTypes()[rq.ServiceTypeId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	if userBalance, err := d.balance.Get(&domain.GetBalanceRequest{UserId: userId}); err == nil {

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
	if _, err := d.balance.Lock(&domain.ModifyBalanceRequest{
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Quantity:      1,
	}); err != nil {
		return nil, err
	}

	// create a delivery object
	delivery := &domain.Delivery{
		Id:            kit.NewId(),
		UserId:        userId,
		ServiceTypeId: rq.ServiceTypeId,
		Status:        "in-progress",
		StartTime:     time.Now().UTC(),
		FinishTime:    nil,
		Details:       rq.Details,
	}

	// save to storage
	delivery, err := d.storage.CreateDelivery(delivery)
	if err != nil {
		return nil, err
	}

	// execute a handler for corresponding task
	if st, ok := d.balance.GetTypes()[rq.ServiceTypeId]; ok {
		if st.DeliveryWfId != "" {
			// start WF
			variables := make(map[string]interface{})
			variables["deliveryId"] = delivery.Id
			_, err := d.bpService.StartProcess(st.DeliveryWfId, variables)
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
	delivery, err = d.storage.UpdateDelivery(delivery)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (d *deliveryServiceImpl) Get(deliveryId string) *domain.Delivery {
	return d.storage.GetDelivery(deliveryId)
}

func (d *deliveryServiceImpl) Complete(deliveryId string, finishTime *time.Time) (*domain.Delivery, error) {

	delivery := d.storage.GetDelivery(deliveryId)
	delivery.Status = "completed"
	if finishTime != nil {
		delivery.FinishTime = finishTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	delivery, err := d.storage.UpdateDelivery(delivery)
	if err != nil {
		return nil, err
	}

	_, err = d.balance.WriteOff(&domain.ModifyBalanceRequest{
		UserId:        delivery.UserId,
		ServiceTypeId: delivery.ServiceTypeId,
		Quantity:      1,
	})
	if err != nil {
		return nil, err
	}

	return delivery, nil

}

func (d *deliveryServiceImpl) Cancel(deliveryId string, cancelTime *time.Time) (*domain.Delivery, error) {

	delivery := d.storage.GetDelivery(deliveryId)
	delivery.Status = "canceled"
	if cancelTime != nil {
		delivery.FinishTime = cancelTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	return d.storage.UpdateDelivery(delivery)
}

func (d *deliveryServiceImpl) UpdateDetails(deliveryId string, details map[string]interface{}) (*domain.Delivery, error) {
	return d.storage.UpdateDetails(deliveryId, details)
}
