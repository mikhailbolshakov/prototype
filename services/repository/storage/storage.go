package storage

import (
	"gitlab.medzdrav.ru/prototype/services/infrastructure"
	"time"
)

type Storage interface {
	CreateBalance(b *Balance) (*Balance, error)
	UpdateBalance(b *Balance) (*Balance, error)
	GetBalance(userId string, at *time.Time) ([]Balance, error)
	GetBalanceForServiceType(userId string, serviceTypeId string, at *time.Time) ([]Balance, error)
	GetTypes() []ServiceType
	CreateDelivery(d *Delivery) (*Delivery, error)
	UpdateDelivery(d *Delivery) (*Delivery, error)
}

type storageImpl struct {
	infr *infrastructure.Container
}

func NewStorage(infr *infrastructure.Container) Storage {
	return &storageImpl{
		infr: infr,
	}
}

func (s *storageImpl) CreateBalance(b *Balance) (*Balance, error) {

	t := time.Now()
	b.CreatedAt, b.UpdatedAt = t, t

	result := s.infr.Db.Instance.Create(b)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) UpdateBalance(b *Balance) (*Balance, error) {

	b.UpdatedAt = time.Now()

	result := s.infr.Db.Instance.Save(b)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) GetBalance(userId string, at *time.Time) ([]Balance, error) {

	var balances []Balance
	result := s.infr.Db.Instance.Where("client_id = ?", userId).Find(&balances)

	if result.Error != nil {
		return nil, result.Error
	}

	return balances, nil

}

func (s *storageImpl) GetBalanceForServiceType(userId string, serviceTypeId string, at *time.Time) ([]Balance, error) {

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

func (s *storageImpl) CreateDelivery(d *Delivery) (*Delivery, error) {

	t := time.Now()
	d.CreatedAt, d.UpdatedAt = t, t

	result := s.infr.Db.Instance.Create(d)

	if result.Error != nil {
		return nil, result.Error
	}

	return d, nil
}

func (s *storageImpl) UpdateDelivery(d *Delivery) (*Delivery, error) {

	d.UpdatedAt = time.Now()

	result := s.infr.Db.Instance.Save(d)

	if result.Error != nil {
		return nil, result.Error
	}

	return d, nil
}