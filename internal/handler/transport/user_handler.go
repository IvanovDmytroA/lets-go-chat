package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IvanovDmytroA/lets-go-chat/internal/service"
)

const minUserNameLength int = 4
const minPasswordLength int = 8
const tokenExpirationHeader string = "X-Expires-After"
const allowedCallsHeader string = "X-Rate-Limit"
const loginCallsLimit string = "5"
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

	userResponse, status := service.CreateUser(userName, password)

	w.Header().Set(contentTypeHeader, jsonType)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&userResponse)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	userName, password := retrieveFormFields(r)
	if len(userName) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userLoginResponse, status := service.LoginUser(userName, password)

	w.WriteHeader(status)
	w.Header().Add(tokenExpirationHeader, time.Now().UTC().String())
	w.Header().Add(allowedCallsHeader, loginCallsLimit)
	json.NewEncoder(w).Encode(&userLoginResponse)
}

func retrieveFormFields(r *http.Request) (userName string, password string) {
	userName = r.FormValue(userNameFormKey)
	password = r.FormValue(passwordFormKey)
	return
}
