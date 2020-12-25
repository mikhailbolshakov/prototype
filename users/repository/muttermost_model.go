package repository

type MMCreateUserRequest struct {
	Username string
	Email string
}

type MMCreateUserResponse struct {
	Id string
}
