package domain

import (
	"errors"
	"fmt"
	"github.com/xtgo/uuid"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/services/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"time"
)

type UserBalanceService interface {
	// get available service types
	GetTypes() map[string]ServiceType
	// adds service to balance
	Add(rq *ModifyBalanceRequest) (*UserBalance, error)
	// requests a balance
	Get(rq *GetBalanceRequest) (*UserBalance, error)
	// write off services
	WriteOff(rq *ModifyBalanceRequest) (*UserBalance, error)
	// lock service
	Lock(rq *ModifyBalanceRequest) (*UserBalance, error)
	// cancel locked service
	Cancel(rq *ModifyBalanceRequest) (*UserBalance, error)
}

func NewBalanceService(userService users.Service, storage storage.Storage, queue queue.Queue) UserBalanceService {
	return &balanceServiceImpl{
		userService: userService,
		storage:     storage,
		queue:       queue,
	}
}

type balanceServiceImpl struct {
	queue       queue.Queue
	storage     storage.Storage
	userService users.Service
}

func (s *balanceServiceImpl) userIdName(input string) string {

	if _, err := uuid.Parse(input); err == nil {
		return input
	} else {
		return s.userService.Get(input).Id
	}

}

func (s *balanceServiceImpl) Add(rq *ModifyBalanceRequest) (*UserBalance, error) {

	userId := s.userIdName(rq.UserId)

	types := s.GetTypes()
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(userId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		_, err = s.storage.CreateBalance(&storage.Balance{
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
		_, err = s.storage.UpdateBalance(&b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(userId, nil)
}

func (s *balanceServiceImpl) get(userId string, at *time.Time) (*UserBalance, error) {

	balanceDtos, err := s.storage.GetBalance(userId, at)
	if err != nil {
		return nil, err
	}

	return s.balanceFromDto(userId, balanceDtos), nil
}

func (s *balanceServiceImpl) Get(rq *GetBalanceRequest) (*UserBalance, error) {
	return s.get(s.userIdName(rq.UserId), nil)
}

func (s *balanceServiceImpl) GetTypes() map[string]ServiceType {

	typesDto := s.storage.GetTypes()
	res := make(map[string]ServiceType, len(typesDto))

	for _, t := range typesDto {
		res[t.Id] = ServiceType{
			Id:           t.Id,
			Description:  t.Description,
			DeliveryWfId: t.DeliveryWfId,
		}
	}
	return res

}

func (s *balanceServiceImpl) WriteOff(rq *ModifyBalanceRequest) (*UserBalance, error) {

	userId := s.userIdName(rq.UserId)

	types := s.GetTypes()
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(userId, rq.ServiceTypeId, nil)
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

		_, err = s.storage.UpdateBalance(&b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(userId, nil)
}

func (s *balanceServiceImpl) Lock(rq *ModifyBalanceRequest) (*UserBalance, error) {

	userId := s.userIdName(rq.UserId)

	types := s.GetTypes()
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(userId, rq.ServiceTypeId, nil)
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
		_, err = s.storage.UpdateBalance(&b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(userId, nil)
}

func (s *balanceServiceImpl) Cancel(rq *ModifyBalanceRequest) (*UserBalance, error) {

	userId := s.userIdName(rq.UserId)

	types := s.GetTypes()
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetBalanceForServiceType(userId, rq.ServiceTypeId, nil)
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
		_, err = s.storage.UpdateBalance(&b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(userId, nil)
}
