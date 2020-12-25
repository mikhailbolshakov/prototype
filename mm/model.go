package mm

type CreateUserRequest struct {
	TeamName string
	Username string
	Password string
	Email string
}

type CreateUserResponse struct {
	Id string
}