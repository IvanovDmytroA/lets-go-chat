package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/IvanovDmytroA/lets-go-chat/internal/service"
	"github.com/labstack/echo"
)

// Create user request handler
func CreateUser(c echo.Context) error {
	ur := new(handler.CreateUserRequest)
	if err := c.Bind(ur); err != nil {
		errMsg := "Error during request body decoding: "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	userResponse, status := service.CreateUser(ur.UserName, ur.Password)

	// Setting up headers
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(status)

	if status != http.StatusOK {
		errMsg := "Failed to create new user"
		log.Printf(errMsg)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	// Encoding and sending response
	enc := json.NewEncoder(c.Response())
	err := enc.Encode(userResponse)
	if err != nil {
		errMsg := "Internal server error. Failed to encode createUserResponse: "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errMsg+err.Error())
	}
	return nil
}
