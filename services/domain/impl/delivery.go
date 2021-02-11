package impl

import (
	"context"
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

func (s *deliveryServiceImpl) userIdName(ctx context.Context, input string) string {
	if _, err := uuid.Parse(input); err == nil {
		return input
	} else {
		return s.userService.Get(ctx, input).Id
	}
}

func (d *deliveryServiceImpl) Delivery(ctx context.Context, rq *domain.DeliveryRequest) (*domain.Delivery, error) {

	userId := d.userIdName(ctx, rq.UserId)

	st, ok := d.balance.GetTypes(ctx)[rq.ServiceTypeId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	if userBalance, err := d.balance.Get(ctx, &domain.GetBalanceRequest{UserId: userId}); err == nil {

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
	if _, err := d.balance.Lock(ctx, &domain.ModifyBalanceRequest{
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
	delivery, err := d.storage.CreateDelivery(ctx, delivery)
	if err != nil {
		return nil, err
	}

	// execute a handler for corresponding task
	if st, ok := d.balance.GetTypes(ctx)[rq.ServiceTypeId]; ok {
		if st.DeliveryWfId != "" {
			// start WF
			variables := make(map[string]interface{})
			variables["deliveryId"] = delivery.Id
			_, err := d.bpService.StartProcess(ctx, st.DeliveryWfId, variables)
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
	delivery, err = d.storage.UpdateDelivery(ctx, delivery)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (d *deliveryServiceImpl) Get(ctx context.Context, deliveryId string) *domain.Delivery {
	return d.storage.GetDelivery(ctx, deliveryId)
}

func (d *deliveryServiceImpl) Complete(ctx context.Context, deliveryId string, finishTime *time.Time) (*domain.Delivery, error) {

	delivery := d.storage.GetDelivery(ctx, deliveryId)
	delivery.Status = "completed"
	if finishTime != nil {
		delivery.FinishTime = finishTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	delivery, err := d.storage.UpdateDelivery(ctx, delivery)
	if err != nil {
		return nil, err
	}

	_, err = d.balance.WriteOff(ctx, &domain.ModifyBalanceRequest{
		UserId:        delivery.UserId,
		ServiceTypeId: delivery.ServiceTypeId,
		Quantity:      1,
	})
	if err != nil {
		return nil, err
	}

	return delivery, nil

}

func (d *deliveryServiceImpl) Cancel(ctx context.Context, deliveryId string, cancelTime *time.Time) (*domain.Delivery, error) {

	delivery := d.storage.GetDelivery(ctx, deliveryId)
	delivery.Status = "canceled"
	if cancelTime != nil {
		delivery.FinishTime = cancelTime
	} else {
		t := time.Now().UTC()
		delivery.FinishTime = &t
	}

	return d.storage.UpdateDelivery(ctx, delivery)
}

func (d *deliveryServiceImpl) UpdateDetails(ctx context.Context, deliveryId string, details map[string]interface{}) (*domain.Delivery, error) {
	return d.storage.UpdateDetails(ctx, deliveryId, details)
}
