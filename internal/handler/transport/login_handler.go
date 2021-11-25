package handler

import (
	"log"
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/IvanovDmytroA/lets-go-chat/internal/service"
	"github.com/labstack/echo"
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

// Login user request handler
func LoginUser(c echo.Context) error {
	ur := new(handler.LoginUserRequest)
	if err := c.Bind(ur); err != nil {
		errMsg := "Error during request body decoding: "
		log.Printf("%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	err := service.LoginUser(ur.UserName, ur.Password, c)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}
