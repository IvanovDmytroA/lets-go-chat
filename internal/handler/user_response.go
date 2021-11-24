package handler

type CreateUserResponse struct {
	Id       string
	UserName string
}

type LoginUserResponse struct {
	Url string
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	Expires     int64
}
