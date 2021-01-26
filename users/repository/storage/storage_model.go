package storage

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
)

type User struct {
	kit.BaseDto
	Id       string `gorm:"column:id"`
	Type     string `gorm:"column:type"`
	Status   string `gorm:"column:status"`
	Username string `gorm:"column:username"`
	MMUserId string `gorm:"column:mm_id"`
	KKUserId string `gorm:"column:kk_id"`
	Details  string `gorm:"column:details"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType    string `gorm:"column:type"`
	Username    string `gorm:"column:username"`
	Status      string `gorm:"column:status"`
	Email       string `gorm:"column:email"`
	Phone       string `gorm:"column:phone"`
	MMId        string `gorm:"column:mm_id"`
	MMChannelId string `gorm:"column:mm_channel_id"`
}

type SearchResponse struct {
	*common.PagingResponse
	Users []*User
}
