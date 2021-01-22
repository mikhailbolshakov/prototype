package users

type CreateUserRequest struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type User struct {
	Id          string `json:"id,omitempty"`
	Type        string `json:"type"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	MMId        string `json:"mmId,omitempty"`
	MMChannelId string `json:"mmChannelId,omitempty"`
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
