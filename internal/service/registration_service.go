package service

import (
	"fmt"
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
	_, exists := userRepo.GetUserByUserName(userName)
	if exists {
		return userResponse, http.StatusBadRequest
	}

	createRequest := handler.CreateUserRequest{UserName: userName, Password: password}
	hash, err := hasher.HashPassword(createRequest.Password)
	if err != nil {
		return userResponse, http.StatusInternalServerError
	}

	userResponse.UserName = createRequest.UserName
	userUuid := uuid.NewV4()
	userResponse.Id = userUuid.String()
	user := model.User{
		Id:       userResponse.Id,
		UserName: userResponse.UserName,
		Password: hash,
	}

	err = userRepo.SaveUser(user)
	if err != nil {
		fmt.Println(err)
		return userResponse, http.StatusBadRequest
	}

	return userResponse, http.StatusOK
}
