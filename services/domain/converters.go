package domain

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/services/repository/storage"
)

func (s *balanceServiceImpl) balanceFromDto(userId string, dtos []storage.Balance) *UserBalance {

	rs := &UserBalance{
		UserId:  userId,
		Balance: map[ServiceType]Balance{},
	}

	types := s.GetTypes()

	for _, d := range dtos {

		rs.Balance[types[d.ServiceTypeId]] = Balance{
			Available: d.Total - d.Delivered - d.Locked,
			Locked:    d.Locked,
			Total:     d.Total,
			Delivered: d.Delivered,
		}

	}

	return rs

}

func (s *deliveryServiceImpl) deliveryToDto(d *Delivery) *storage.Delivery {

	detailsJson := "{}"
	if d.Details != nil {
		bytes, _ := json.Marshal(d.Details)
		detailsJson = string(bytes)
	}

	return &storage.Delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     d.StartTime,
		FinishTime:    d.FinishTime,
		Details:       detailsJson,
	}
}

func (s *deliveryServiceImpl) deliveryFromDto(d *storage.Delivery) *Delivery {

	var v map[string]interface{}
	if d.Details != "" {
		_ = json.Unmarshal([]byte(d.Details), &v)
	}

	return &Delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     d.StartTime,
		FinishTime:    d.FinishTime,
		Details:       v,
	}
}