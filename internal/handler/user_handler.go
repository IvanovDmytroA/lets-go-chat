package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/IvanovDmytroA/lets-go-chat/pkg/hasher"
	uuid "github.com/nu7hatch/gouuid"
)

const minUserNameLength int = 4
const minPasswordLength int = 8
const tokenExpirationHeader string = "X-Expires-After"
const allowedCallsHeader string = "X-Rate-Limit"
const loginCallsLimit string = "5"
const url string = "Link to the chat"
const contentTypeHeader string = "Content-Type"
const jsonType string = "application/json"
const userNameFormKey string = "userName"
const passwordFormKey string = "password"

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userName, password := retrieveFormFields(r)
	if len(userName) < minUserNameLength || len(password) < minPasswordLength {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := repository.GetUserByUserName(userName)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userUuid, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createRequest := model.CreateUserRequest{UserName: userName, Password: password}
	hash, err := hasher.HashPassword(createRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResponse := model.CreateUserResponse{
		UserName: createRequest.UserName,
		Id:       userUuid.String(),
	}
	user := model.User{
		Id:       userResponse.Id,
		UserName: userResponse.UserName,
		Password: hash,
	}

	repository.SaveUser(user)

	w.Header().Set(contentTypeHeader, jsonType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&userResponse)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	userName, password := retrieveFormFields(r)
	if len(userName) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginRequest := model.LoginUserRequest{UserName: userName, Password: password}
	responseStatus := defineResponseStatus(loginRequest)
	w.WriteHeader(responseStatus)
	if responseStatus != http.StatusOK {
		return
	}

	w.Header().Add(tokenExpirationHeader, time.Now().UTC().String())
	w.Header().Add(allowedCallsHeader, loginCallsLimit)

	userLoginResponse := model.LoginUserResponse{Url: url}

	json.NewEncoder(w).Encode(&userLoginResponse)
}

func defineResponseStatus(loginRequest model.LoginUserRequest) int {
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

func retrieveFormFields(r *http.Request) (userName string, password string) {
	userName = r.FormValue(userNameFormKey)
	password = r.FormValue(passwordFormKey)
	return
}
