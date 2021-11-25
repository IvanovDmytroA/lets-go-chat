package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/labstack/echo"
)

func GetActiveUsers(c echo.Context) error {
	activeUsers := repository.GetActiveUsersStorage().Users
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	json.NewEncoder(c.Response()).Encode(&activeUsers)
	return nil
}
