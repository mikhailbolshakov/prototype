package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/xtgo/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"time"
)

func NewBalanceService(userService domain.UserService, storage domain.Storage, queue queue.Queue) domain.UserBalanceService {
	return &balanceServiceImpl{
		userService: userService,
		storage:     storage,
		queue:       queue,
	}
}

type balanceServiceImpl struct {
	queue       queue.Queue
	storage     domain.Storage
	userService domain.UserService
}

func (s *balanceServiceImpl) userIdName(ctx context.Context, input string) string {

	if _, err := uuid.Parse(input); err == nil {
		return input
	} else {
		return s.userService.Get(ctx, input).Id
	}

}

func (s *balanceServiceImpl) Add(ctx context.Context, rq *domain.ModifyBalanceRequest) (*domain.UserBalance, error) {

	userId := s.userIdName(ctx, rq.UserId)

	types := s.GetTypes(ctx)
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(ctx, userId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		_, err = s.storage.CreateBalance(ctx, &domain.BalanceItem{
			Id:            kit.NewId(),
			UserId:        userId,
			ServiceTypeId: rq.ServiceTypeId,
			Total:         rq.Quantity,
		})
		if err != nil {
			return nil, err
		}
	} else if len(balances) == 1 {
		b := balances[0]
		b.Total = b.Total + rq.Quantity
		_, err = s.storage.UpdateBalance(ctx, b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(ctx, userId, nil)
}

func (s *balanceServiceImpl) toUserBalance(ctx context.Context, userId string, items []*domain.BalanceItem) *domain.UserBalance {

	rs := &domain.UserBalance{
		UserId:  userId,
		Balance: map[domain.ServiceType]domain.Balance{},
	}

	types := s.GetTypes(ctx)

	for _, d := range items {
		rs.Balance[types[d.ServiceTypeId]] = domain.Balance{
			Available: d.Total - d.Delivered - d.Locked,
			Locked:    d.Locked,
			Total:     d.Total,
			Delivered: d.Delivered,
		}

	}

	return rs

}

func (s *balanceServiceImpl) get(ctx context.Context, userId string, at *time.Time) (*domain.UserBalance, error) {

	balanceItems, err := s.storage.GetBalance(ctx, userId, at)
	if err != nil {
		return nil, err
	}

	return s.toUserBalance(ctx, userId, balanceItems), nil
}

func (s *balanceServiceImpl) Get(ctx context.Context, rq *domain.GetBalanceRequest) (*domain.UserBalance, error) {
	return s.get(ctx, s.userIdName(ctx, rq.UserId), nil)
}

func (s *balanceServiceImpl) GetTypes(ctx context.Context) map[string]domain.ServiceType {

	typesDto := s.storage.GetTypes(ctx)
	res := make(map[string]domain.ServiceType, len(typesDto))

	for _, t := range typesDto {
		res[t.Id] = domain.ServiceType{
			Id:           t.Id,
			Description:  t.Description,
			DeliveryWfId: t.DeliveryWfId,
		}
	}
	return res

}

func (s *balanceServiceImpl) WriteOff(ctx context.Context, rq *domain.ModifyBalanceRequest) (*domain.UserBalance, error) {

	userId := s.userIdName(ctx, rq.UserId)

	types := s.GetTypes(ctx)
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(ctx, userId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return nil, errors.New(fmt.Sprintf("cannot write off service %s, no availables", rq.ServiceTypeId))
	} else if len(balances) == 1 {

		b := balances[0]

		if b.Locked < rq.Quantity {
			return nil, errors.New(fmt.Sprintf("only locked service can be written of"))
		}

		b.Locked = b.Locked - rq.Quantity
		b.Delivered = b.Delivered + rq.Quantity

		_, err = s.storage.UpdateBalance(ctx, b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(ctx, userId, nil)
}

func (s *balanceServiceImpl) Lock(ctx context.Context, rq *domain.ModifyBalanceRequest) (*domain.UserBalance, error) {

	userId := s.userIdName(ctx, rq.UserId)

	types := s.GetTypes(ctx)
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(ctx, userId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return nil, errors.New(fmt.Sprintf("cannot lock service %s, no availables", rq.ServiceTypeId))
	} else if len(balances) == 1 {
		b := balances[0]

		if rq.Quantity > (b.Total - b.Locked) {
			return nil, errors.New("cannot lock more then availables")
		}

		b.Locked = b.Locked + rq.Quantity
		_, err = s.storage.UpdateBalance(ctx, b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(ctx, userId, nil)
}

func (s *balanceServiceImpl) Cancel(ctx context.Context, rq *domain.ModifyBalanceRequest) (*domain.UserBalance, error) {

	userId := s.userIdName(ctx, rq.UserId)

	types := s.GetTypes(ctx)
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(ctx, userId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return nil, errors.New(fmt.Sprintf("cannot cancel service %s, no locked", rq.ServiceTypeId))
	} else if len(balances) == 1 {
		b := balances[0]

		if rq.Quantity > b.Locked {
			return nil, errors.New("cannot cancel more then locked")
		}

		b.Locked = b.Locked - rq.Quantity
		_, err = s.storage.UpdateBalance(ctx, b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(ctx, userId, nil)
}
