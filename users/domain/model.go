package domain

import "gitlab.medzdrav.ru/prototype/kit/common"

const (
	USER_TYPE_CLIENT     = "client"
	USER_TYPE_CONSULTANT = "consultant"
	USER_TYPE_EXPERT     = "expert"
	USER_TYPE_SUPERVISOR = "supervisor"
)

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Type        string `json:"type"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	MMUserId    string `json:"mmId"`
	MMChannelId string `json:"mmChannelId"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType       string
	Username       string
	Email          string
	Phone          string
	MMId           string
	MMChannelId    string
	OnlineStatuses []string
}

type SearchResponse struct {
	*common.PagingResponse
	Users []*User
}
