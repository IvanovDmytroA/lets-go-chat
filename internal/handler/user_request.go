package handler

type CreateUserRequest struct {
	UserName string
	Password string
}

type LoginUserRequest struct {
	UserName string
	Password string
}
