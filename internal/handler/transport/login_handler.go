package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

	service.LoginUser(ur.UserName, ur.Password, c)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().Header().Set("X-Rate-Limit", strconv.Itoa(360))
	c.Response().Header().Set("X-Expires-After", time.Now().Add(time.Minute*10).Format(time.RFC1123))
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
