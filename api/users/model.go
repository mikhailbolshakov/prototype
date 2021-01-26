package users

import "time"

type CreateClientRequest struct {
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Sex        string    `json:"sex"`
	BirthDate  time.Time `json:"birthDate"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
}

type CreateConsultantRequest struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
}

type CreateExpertRequest struct {
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
	MMUserId          string             `json:"mmId,omitempty"`
	KKUserId          string             `json:"kkId,omitempty"`
	ClientDetails     *ClientDetails     `json:"clientDetails,omitempty"`
	ConsultantDetails *ConsultantDetails `json:"consultantDetails,omitempty"`
	ExpertDetails     *ExpertDetails     `json:"expertDetails,omitempty"`
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
	MMChannelId       string             `json:"mmChannelId,omitempty"`
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
	Specialization string `json:"specialization,omitempty"`
}

type SearchResponse struct {
	Index int     `json:"index"`
	Total int     `json:"total"`
	Users []*User `json:"users"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AuthCode string `json:"authCode"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}
