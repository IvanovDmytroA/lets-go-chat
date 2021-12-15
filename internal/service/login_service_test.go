package service

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	td "github.com/IvanovDmytroA/lets-go-chat/tests"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/uptrace/bun"
)

var userUuid = uuid.NewV4()
var user = model.User{
	Id:       userUuid.String(),
	UserName: "user",
	Password: "$2a$14$NuOsuc8HUa1tvdcQPQQ4H.zlobVv3S/oVe/ln5G69en0D7a8DAqLi",
}

func TestLoginUser(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	td.DelUser(dbConnect, t)

	err := createTestUser(dbConnect)
	if err != nil {
		t.Fatalf("Saving test user failed")
	}

	e := echo.New()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user/login", bytes.NewReader(td.TestDataM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	errLogin := LoginUser("user", "password", c)

	if errLogin != nil {
		t.Fatalf("Failed to login user")
	}

	td.DelUser(dbConnect, t)

}

func createTestUser(dbc *bun.DB) error {
	userRep := repository.GetUsersRepo()
	err := userRep.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}
