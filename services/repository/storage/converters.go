package storage

import (
	"encoding/json"
	"gitlab.medzdrav.ru/prototype/services/domain"
)

func (s *storageImpl) toBalanceItemDto(bi *domain.BalanceItem) *balanceItem {
	if bi == nil {
		return nil
	}
	return &balanceItem{
		Id:            bi.Id,
		UserId:        bi.UserId,
		ServiceTypeId: bi.ServiceTypeId,
		Total:         bi.Total,
		Delivered:     bi.Delivered,
		Locked:        bi.Locked,
	}
}

func (s *storageImpl) toBalanceItemDomain(bi *balanceItem) *domain.BalanceItem {
	if bi == nil {
		return nil
	}
	return &domain.BalanceItem {
		Id:            bi.Id,
		UserId:        bi.UserId,
		ServiceTypeId: bi.ServiceTypeId,
		Total:         bi.Total,
		Delivered:     bi.Delivered,
		Locked:        bi.Locked,
	}
}

func (s *storageImpl) toBalanceItemsDomain(dtos []*balanceItem) []*domain.BalanceItem {
	var res []*domain.BalanceItem
	for _, d := range dtos {
		res = append(res, s.toBalanceItemDomain(d))
	}
	return res
}

func (s *storageImpl) toServiceTypeDomain(dto serviceType) domain.ServiceType {
	return domain.ServiceType{
		Id:           dto.Id,
		Description:  dto.Description,
		DeliveryWfId: dto.DeliveryWfId,
	}
}

func (s *storageImpl) toServiceTypesDomain(dtos []serviceType) []domain.ServiceType {
	var res []domain.ServiceType
	for _, d := range dtos {
		res = append(res, s.toServiceTypeDomain(d))
	}
	return res
}

func (s *storageImpl) toDeliveryDto(d *domain.Delivery) *delivery {

	detailsJson := "{}"
	if d.Details != nil {
		bytes, _ := json.Marshal(d.Details)
		detailsJson = string(bytes)
	}

	return &delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     d.StartTime,
		FinishTime:    d.FinishTime,
		Details:       detailsJson,
	}
}

func (s *storageImpl) toDeliveryDomain(d *delivery) *domain.Delivery {

	var v map[string]interface{}
	if d.Details != "" {
		_ = json.Unmarshal([]byte(d.Details), &v)
	}

	return &domain.Delivery{
		Id:            d.Id,
		UserId:        d.UserId,
		ServiceTypeId: d.ServiceTypeId,
		Status:        d.Status,
		StartTime:     d.StartTime,
		FinishTime:    d.FinishTime,
		Details:       v,
	}
}