package service

import (
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/IvanovDmytroA/lets-go-chat/pkg/hasher"
	uuid "github.com/satori/go.uuid"
)

func CreateUser(userName, password string) (handler.CreateUserResponse, int) {
	userResponse := handler.CreateUserResponse{}
	userRepo := repository.GetUsersRepo()
	_, err := userRepo.GetUserByUserName(userName)
	if err != nil {
		return userResponse, http.StatusBadRequest
	}

	userUuid := uuid.NewV4()
	if err != nil {
		return userResponse, http.StatusInternalServerError
	}

	createRequest := handler.CreateUserRequest{UserName: userName, Password: password}
	hash, err := hasher.HashPassword(createRequest.Password)
	if err != nil {
		return userResponse, http.StatusInternalServerError
	}

	userResponse.UserName = createRequest.UserName
	userResponse.Id = userUuid.String()
	user := model.User{
		Id:       userResponse.Id,
		UserName: userResponse.UserName,
		Password: hash,
	}

	userRepo.SaveUser(user)

	return userResponse, http.StatusOK
}
