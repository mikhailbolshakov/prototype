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
	Id        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	MMId      string `json:"mmId,omitempty"`
}


