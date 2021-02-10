package storage

import (
	"encoding/json"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"gitlab.medzdrav.ru/prototype/services/domain"
	"time"
)


type serviceType struct {
	Id           string `gorm:"column:id"`
	Description  string `gorm:"column:description"`
	DeliveryWfId string `gorm:"column:delivery_wf_id"`
}

type balanceItem struct {
	kit.BaseDto
	Id            string `gorm:"column:id"`
	UserId        string `gorm:"column:client_id"`
	ServiceTypeId string `gorm:"column:service_type_id"`
	Total         int    `gorm:"column:total"`
	Delivered     int    `gorm:"column:delivered"`
	Locked        int    `gorm:"column:locked"`
}

type delivery struct {
	kit.BaseDto
	Id            string     `gorm:"column:id"`
	UserId        string     `gorm:"column:client_id"`
	ServiceTypeId string     `gorm:"column:service_type_id"`
	Status        string     `gorm:"column:status"`
	StartTime     time.Time  `gorm:"column:start_time"`
	FinishTime    *time.Time `gorm:"column:finish_time"`
	Details       string     `gorm:"column:details"`
}

type storageImpl struct {
	c *container
}

func newStorage(c *container) *storageImpl {
	return &storageImpl{c}
}

func (s *storageImpl) CreateBalance(b *domain.BalanceItem) (*domain.BalanceItem, error) {

	dto := s.toBalanceItemDto(b)

	t := time.Now().UTC()
	dto.CreatedAt, dto.UpdatedAt = t, t

	result := s.c.Db.Instance.Create(dto)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) UpdateBalance(b *domain.BalanceItem) (*domain.BalanceItem, error) {

	dto := s.toBalanceItemDto(b)

	dto.UpdatedAt = time.Now().UTC()

	result := s.c.Db.Instance.Save(dto)

	if result.Error != nil {
		return nil, result.Error
	}

	return b, nil

}

func (s *storageImpl) GetBalance(userId string, at *time.Time) ([]*domain.BalanceItem, error) {

	var balances []*balanceItem
	result := s.c.Db.Instance.Where("client_id = ?", userId).Find(&balances)

	if result.Error != nil {
		return nil, result.Error
	}

	return s.toBalanceItemsDomain(balances), nil

}

func (s *storageImpl) GetBalanceForServiceType(userId string, serviceTypeId string, at *time.Time) ([]*domain.BalanceItem, error) {

	var balances []*balanceItem
	result := s.c.Db.Instance.
		Where("client_id = ?", userId).
		Where("service_type_id = ?", serviceTypeId).
		Find(&balances)

	if result.Error != nil {
		return nil, result.Error
	}

	return s.toBalanceItemsDomain(balances), nil
}

func (s *storageImpl) GetTypes() []domain.ServiceType {
	// TODO: cache
	var types []serviceType
	s.c.Db.Instance.Find(&types)
	return s.toServiceTypesDomain(types)
}

func (s *storageImpl) CreateDelivery(d *domain.Delivery) (*domain.Delivery, error) {
	dto := s.toDeliveryDto(d)
	t := time.Now().UTC()
	dto.CreatedAt, dto.UpdatedAt = t, t
	result := s.c.Db.Instance.Create(dto)
	if result.Error != nil {
		return nil, result.Error
	}
	return d, nil
}

func (s *storageImpl) UpdateDelivery(d *domain.Delivery) (*domain.Delivery, error) {
	dto := s.toDeliveryDto(d)
	dto.UpdatedAt = time.Now().UTC()
	result := s.c.Db.Instance.Save(dto)
	if result.Error != nil {
		return nil, result.Error
	}
	return d, nil
}

func (s *storageImpl) UpdateDetails(deliveryId string, details map[string]interface{}) (*domain.Delivery, error) {

	dt := "{}"
	if details != nil {
		bytes, _ := json.Marshal(details)
		dt = string(bytes)
	}

	d := &delivery{	Id: deliveryId}
	result := s.c.Db.Instance.Model(d).Updates(map[string]interface{}{"details": dt, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return nil, result.Error
	}
	return s.GetDelivery(deliveryId), nil
}

func (s *storageImpl) GetDelivery(id string) *domain.Delivery {
	res := &delivery{Id: id}
	s.c.Db.Instance.First(res)
	return s.toDeliveryDomain(res)
}
