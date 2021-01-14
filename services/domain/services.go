package domain

import (
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
	"time"
)

type UserBalanceService interface {
	// adds service to balance
	Add(rq *ModifyBalanceRequest) (*UserBalance, error)
	// requests a balance
	Get(rq *GetBalanceRequest) (*UserBalance, error)
	// write off services
	WriteOff(rq *ModifyBalanceRequest) (*UserBalance, error)
	// lock service
	Lock(rq *ModifyBalanceRequest) (*UserBalance, error)
}

type DeliveryService interface{}

func NewService(storage storage.BalanceStorage, queue queue.Queue) UserBalanceService {
	return &serviceImpl{
		storage: storage,
		queue:   queue,
	}
}

type serviceImpl struct {
	queue   queue.Queue
	storage storage.BalanceStorage
}

func (s *serviceImpl) Add(rq *ModifyBalanceRequest) (*UserBalance, error) {

	types := s.GetTypes()
	if _, ok := types[rq.ServiceTypeId]; !ok {
		return nil, errors.New(fmt.Sprintf("service type %s isn't supported", rq.ServiceTypeId))
	}

	balances, err := s.storage.GetForServiceType(rq.UserId, rq.ServiceTypeId, nil)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		_, err = s.storage.Create(&storage.Balance{
			Id:            kit.NewId(),
			UserId:        rq.UserId,
			ServiceTypeId: rq.ServiceTypeId,
			Total:         rq.Quantity,
		})
		if err != nil {
			return nil, err
		}
	} else if len(balances) == 1 {
		b := balances[0]
		b.Total = b.Total + rq.Quantity
		_, err = s.storage.Update(&b)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("balance is corrupted")
	}

	return s.get(rq.UserId, nil)
}

func (s *serviceImpl) get(userId string, at *time.Time) (*UserBalance, error) {

	balanceDtos, err := s.storage.Get(userId, at)
	if err != nil {
		return nil, err
	}

	return s.balanceFromDto(userId, balanceDtos), nil
}

func (s *serviceImpl) Get(rq *GetBalanceRequest) (*UserBalance, error) {
	return s.get(rq.UserId, nil)
}

func (s *serviceImpl) GetTypes() map[string]ServiceType {

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

func (s *serviceImpl) WriteOff(rq *ModifyBalanceRequest) (*UserBalance, error) {
	return nil, nil
}
func (s *serviceImpl) Lock(rq *ModifyBalanceRequest) (*UserBalance, error) {
	return nil, nil
}
