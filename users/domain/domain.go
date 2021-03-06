package domain

import (
	"context"
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
	PhotoUrl          string             `json:"photoUrl"`
}

type ConsultantDetails struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	PhotoUrl   string `json:"photoUrl"`
}

type ExpertDetails struct {
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PhotoUrl       string `json:"photoUrl"`
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

type UserService interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByUsername(ctx context.Context, username string) *User
	GetByMMId(ctx context.Context, mmId string) *User
	Get(ctx context.Context, id string) *User
	Activate(ctx context.Context, userId string) (*User, error)
	Delete(ctx context.Context, userId string) (*User, error)
	SetClientDetails(ctx context.Context, userId string, details *ClientDetails) (*User, error)
	SetMMUserId(ctx context.Context, userId, mmId string) (*User, error)
	SetKKUserId(ctx context.Context, userId, kkId string) (*User, error)
	Search(ctx context.Context, cr *SearchCriteria) (*SearchResponse, error)
}

