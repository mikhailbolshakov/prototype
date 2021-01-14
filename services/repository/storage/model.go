package storage

import (
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
	"time"
)

type ServiceType struct {
	Id           string `gorm:"column:id"`
	Description  string `gorm:"column:description"`
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

type Delivery struct {
	kit.BaseDto
	Id            string                 `gorm:"column:id"`
	UserId        string                 `gorm:"column:client_id"`
	ServiceTypeId string                 `gorm:"column:service_type_id"`
	Status        string                 `gorm:"column:status"`
	StartTime     time.Time              `gorm:"column:start_time"`
	FinishTime    *time.Time             `gorm:"column:finish_time"`
	Details       string `gorm:"column:details"`
}
