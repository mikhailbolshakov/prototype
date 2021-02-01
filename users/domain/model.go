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

	USER_GRP_CLIENT             = "client"
	USER_GRP_CONSULTANT_LAWYER  = "consultant-lawyer"
	USER_GRP_CONSULTANT_DOCTOR  = "medconsultant"
	USER_GRP_DOCTOR_DENTIST     = "doctor-dentist"
	USER_GRP_DOCTOR_PRTHOPEDIST = "doctor-orthopedist"
	USER_GRP_SUPERVIZOR_RGS     = "supervisor-rgs"
)

const (
	USER_STATUS_DRAFT   = "draft"
	USER_STATUS_ACTIVE  = "active"
	USER_STATUS_LOCKED  = "locked"
	USER_STATUS_DELETED = "deleted"
)

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
	CommonChannelId   string             `json:"commonChannelId"`
	MedChannelId      string             `json:"medChannelId"`
	LawChannelId      string             `json:"lawChannelId"`
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

type User struct {
	Id                string             `json:"id"`
	Username          string             `json:"username"`
	Type              string             `json:"type"`
	Status            string             `json:"status"`
	MMUserId          string             `json:"mmId"`
	KKUserId          string             `json:"kkId"`
	ClientDetails     *ClientDetails     `json:"clientDetails,omitempty"`
	ConsultantDetails *ConsultantDetails `json:"consultantDetails,omitempty"`
	ExpertDetails     *ExpertDetails     `json:"expertDetails,omitempty"`
	Groups            []string           `json:"groups"`
	CreatedAt         time.Time          `json:"createdAt"`
	ModifiedAt        time.Time          `json:"modifiedAt"`
	DeletedAt         *time.Time         `json:"deletedAt,omitempty"`
}

type SearchCriteria struct {
	*common.PagingRequest
	UserType        string
	Username        string
	UserGroup       string
	Status          string
	Email           string
	Phone           string
	MMId            string
	CommonChannelId string
	MedChannelId    string
	LawChannelId    string
	OnlineStatuses  []string
}

type SearchResponse struct {
	*common.PagingResponse
	Users []*User
}
