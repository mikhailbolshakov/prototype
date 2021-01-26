package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"time"
)

const (
	USER_TYPE_CLIENT     = "client"
	USER_TYPE_CONSULTANT = "consultant"
	USER_TYPE_EXPERT     = "expert"
	USER_TYPE_SUPERVISOR = "supervisor"
)

const (
	USER_STATUS_DRAFT   = "draft"
	USER_STATUS_ACTIVE  = "active"
	USER_STATUS_LOCKED  = "locked"
	USER_STATUS_DELETED = "deleted"
)

var statusMap = map[string]struct{}{
	USER_STATUS_DRAFT:  {},
	USER_STATUS_ACTIVE: {},
	USER_STATUS_LOCKED: {},
}

type User struct {
	Id                string             `json:"id"`
	Username          string             `json:"username"`
	Type              string             `json:"type"`
	Status            string             `json:"status"`
	MMUserId          string             `json:"mmId"`
	KKUserId          string             `json:"kkId"`
	ClientDetails     *ClientDetails     `json:"clientDetails"`
	ConsultantDetails *ConsultantDetails `json:"consultantDetails"`
	ExpertDetails     *ExpertDetails     `json:"expertDetails"`
	ModifiedAt        time.Time          `json:"modifiedAt"`
	DeletedAt         *time.Time         `json:"deletedAt,omitempty"`
}

type PersonalAgreement struct {
	GivenAt   *time.Time `json:"givenAt"`
	RevokedAt *time.Time `json:"revokedAt"`
}

type ClientDetails struct {
	FirstName         string             `json:"firstName"`
	MiddleName        string             `json:"middleName"`
	LastName          string             `json:"lastName"`
	Sex               string             `json:"sex"`
	BirthDate         time.Time          `json:"birthDate"`
	Phone             string             `json:"phone"`
	Email             string             `json:"email"`
	PersonalAgreement *PersonalAgreement `json:"personalAgreement"`
	MMChannelId       string             `json:"mmChannelId"`
}

type ConsultantDetails struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
}

type ExpertDetails struct {
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType       string
	Username       string
	Status         string
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
