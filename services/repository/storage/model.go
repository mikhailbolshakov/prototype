package storage

import kit "gitlab.medzdrav.ru/prototype/kit/storage"

type ServiceType struct {
	Id string `gorm:"column:id"`
	Description string `gorm:"column:description"`
	DeliveryWfId string `gorm:"column:delivery_wf_id"`
}

type Balance struct {
	kit.BaseDto
	Id            string `gorm:"column:id"`
	UserId        string `gorm:"column:client_id"`
	ServiceTypeId string `gorm:"column:service_type_id"`
	Total         int    `gorm:"column:total"`
	Delivered     int    `gorm:"column:delivered"`
	Locked        int    `gorm:"column:locked"`
}
