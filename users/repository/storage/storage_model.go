package storage

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	kit "gitlab.medzdrav.ru/prototype/kit/storage"
)

type UserGroup struct {
	kit.BaseDto
	Id     string `gorm:"column:id"`
	UserId string `gorm:"column:user_id"`
	Type   string `gorm:"column:type"`
	Group  string `gorm:"column:group_code"`
}

type User struct {
	kit.BaseDto
	Id       string       `gorm:"column:id"`
	Type     string       `gorm:"column:type"`
	Status   string       `gorm:"column:status"`
	Username string       `gorm:"column:username"`
	MMUserId string       `gorm:"column:mm_id"`
	KKUserId string       `gorm:"column:kk_id"`
	Details  string       `gorm:"column:details"`
	Groups   []*UserGroup `gorm:"-"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType        string
	UserGroup       string
	Username        string
	Status          string
	Email           string
	Phone           string
	MMId            string
	CommonChannelId string
	MedChannelId    string
	LawChannelId    string
}

type SearchResponse struct {
	*common.PagingResponse
	Users []*User
}
