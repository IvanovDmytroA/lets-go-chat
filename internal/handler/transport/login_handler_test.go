package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	td "github.com/IvanovDmytroA/lets-go-chat/tests"
	"github.com/labstack/echo"
)

func TestLoginUser(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	td.DelUser(dbConnect, t)
	err := td.AddTestUser(dbConnect)
	defer td.DelUser(dbConnect, t)
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

	err = LoginUser(c)

	if err != nil {
		t.Fatalf("Failed to login user: %s", err.Error())
	}
}

func TestLoginUserBadRequestError(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	td.DelUser(dbConnect, t)

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user/login", bytes.NewReader(td.TestDataM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	err := LoginUser(c)

	if err == nil {
		t.Fatal("Expected to get error, but user successfully logged in")
	}
}

func TestLoginUserFailedBindData(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	td.DelUser(dbConnect, t)

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user/login", bytes.NewReader([]byte{}))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	err := LoginUser(c)

	if err == nil {
		t.Fatal("Expected to get error, but user successfully logged in")
	}
}
