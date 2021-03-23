package users

import "time"

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ChatLogin bool   `json:"chatLogin"`
}

type LoginResponse struct {
	SessionId string `json:"sessionId"`
}

type CreateClientRequest struct {
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Sex        string    `json:"sex"`
	BirthDate  time.Time `json:"birthDate"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	PhotoUrl   string    `json:"photoUrl"`
}

type CreateConsultantRequest struct {
	FirstName  string   `json:"firstName"`
	MiddleName string   `json:"middleName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	PhotoUrl   string   `json:"photoUrl"`
	Groups     []string `json:"groups"`
}

type CreateExpertRequest struct {
	FirstName  string   `json:"firstName"`
	MiddleName string   `json:"middleName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	PhotoUrl   string   `json:"photoUrl"`
	Groups     []string `json:"groups"`
}

type User struct {
	Id                string             `json:"id"`
	Username          string             `json:"username"`
	Type              string             `json:"type"`
	Status            string             `json:"status"`
	MMUserId          string             `json:"mmId,omitempty"`
	KKUserId          string             `json:"kkId,omitempty"`
	ClientDetails     *ClientDetails     `json:"clientDetails,omitempty"`
	ConsultantDetails *ConsultantDetails `json:"consultantDetails,omitempty"`
	ExpertDetails     *ExpertDetails     `json:"expertDetails,omitempty"`
	Groups            []string           `json:"groups,omitempty"`
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
	CommonChannelId   string             `json:"commonChannelId,omitempty"`
	MedChannelId      string             `json:"medChannelId,omitempty"`
	LawChannelId      string             `json:"lawChannelId,omitempty"`
	PhotoUrl          string             `json:"photoUrl,omitempty"`
}

type ConsultantDetails struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	PhotoUrl   string `json:"photoUrl,omitempty"`
}

type ExpertDetails struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	PhotoUrl   string `json:"photoUrl,omitempty"`
}

type SearchResponse struct {
	Index int     `json:"index"`
	Total int     `json:"total"`
	Users []*User `json:"users"`
}
