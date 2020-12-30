package storage

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
)

type User struct {
	kit.BaseDto
	Id          string `gorm:"column:id"`
	Type        string `gorm:"column:type"`
	Username    string `gorm:"column:username"`
	FirstName   string `gorm:"column:first_name"`
	LastName    string `gorm:"column:last_name"`
	Phone       string `gorm:"column:phone"`
	Email       string `gorm:"column:email"`
	MMUserId    string `gorm:"column:mm_id"`
	MMChannelId string `gorm:"column:mm_channel_id"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType    string `gorm:"column:type"`
	Username    string `gorm:"column:username"`
	Email       string `gorm:"column:email"`
	Phone       string `gorm:"column:phone"`
	MMId        string `gorm:"column:mm_id"`
	MMChannelId string `gorm:"column:mm_channel_id"`
}

type SearchResponse struct {
	*common.PagingResponse
	Users []*User
}
