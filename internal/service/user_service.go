package service

import (
	"errors"
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/IvanovDmytroA/lets-go-chat/pkg/hasher"
	uuid "github.com/nu7hatch/gouuid"
)

const url string = "Link to the chat"

func CreateUser(userName, password string) (handler.CreateUserResponse, int) {
	userResponse := handler.CreateUserResponse{}
	_, err := repository.GetUserByUserName(userName)
	if err == nil {
		return userResponse, http.StatusBadRequest
	}

	userUuid, err := uuid.NewV4()
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

	repository.SaveUser(user)

	return userResponse, http.StatusOK
}

func LoginUser(userName, password string) (handler.LoginUserResponse, int) {
	userLoginResponse := handler.LoginUserResponse{}
	loginRequest := handler.LoginUserRequest{UserName: userName, Password: password}
	responseStatus := defineResponseStatus(loginRequest)
	if responseStatus != http.StatusOK {
		return userLoginResponse, responseStatus
	}
	userLoginResponse.Url = url
	return userLoginResponse, responseStatus
}

func defineResponseStatus(loginRequest handler.LoginUserRequest) int {
	user, err := repository.GetUserByUserName(loginRequest.UserName)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return http.StatusBadRequest
		} else {
			return http.StatusInternalServerError
		}
	}

	isCorrectPassword := hasher.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isCorrectPassword {
		return http.StatusBadRequest
	}

	return http.StatusOK
}
