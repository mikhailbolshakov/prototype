package storage

import (
	"gitlab.medzdrav.ru/prototype/services/infrastructure"
	"time"
)

type BalanceStorage interface {
	Create(b *Balance) (*Balance, error)
	Update(b *Balance) (*Balance, error)
	Get(userId string, at *time.Time) ([]Balance, error)
	GetForServiceType(userId string, serviceTypeId string, at *time.Time) ([]Balance, error)
	GetTypes() []ServiceType
}

type storageImpl struct {
	infr *infrastructure.Container
}

func NewStorage(infr *infrastructure.Container) BalanceStorage {
	return &storageImpl{
		infr: infr,
	}
}

func (s *storageImpl) Create(b *Balance) (*Balance, error) {

	t := time.Now()
	b.CreatedAt, b.UpdatedAt = t, t

	result := s.infr.Db.Instance.Create(b)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) Update(b *Balance) (*Balance, error) {

	b.UpdatedAt = time.Now()

	result := s.infr.Db.Instance.Save(b)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) Get(userId string, at *time.Time) ([]Balance, error) {

	var balances []Balance
	result := s.infr.Db.Instance.Where("client_id = ?", userId).Find(&balances)

	if result.Error != nil {
		return nil, result.Error
	}

	return balances, nil

}

func (s *storageImpl) GetForServiceType(userId string, serviceTypeId string, at *time.Time) ([]Balance, error) {

	var balances []Balance
	result := s.infr.Db.Instance.
		Where("client_id = ?", userId).
		Where("service_type_id = ?", serviceTypeId).
		Find(&balances)

	if result.Error != nil {
		return nil, result.Error
	}

	return balances, nil
}

func (s *storageImpl) GetTypes() []ServiceType {

	var types []ServiceType
	s.infr.Db.Instance.Find(&types)
	return types

}