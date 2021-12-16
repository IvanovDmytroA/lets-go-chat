package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	td "github.com/IvanovDmytroA/lets-go-chat/tests"
	"github.com/labstack/echo"
)

func TestCreateUser(t *testing.T) {
	dbConnect := td.DBConnection()
	td.DelUser(dbConnect, t)
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user", bytes.NewReader(td.TestDataM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)

	defer td.DelUser(dbConnect, t)

	err := CreateUser(c)

	if err != nil {
		t.Fatalf("Failed to create a user: %s", err.Error())
	}
}

func TestCreateUserFailedBindData(t *testing.T) {
	dbConnect := td.DBConnection()
	td.DelUser(dbConnect, t)
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user", bytes.NewReader([]byte{}))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)

	defer td.DelUser(dbConnect, t)

	err := CreateUser(c)

	if err == nil {
		t.Fatal("Expected error, but user was created")
	}
}

func TestCreateUserFailedIncompleteData(t *testing.T) {
	dbConnect := td.DBConnection()
	td.DelUser(dbConnect, t)
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user", bytes.NewReader(td.IncompleteTestDataM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)

	defer td.DelUser(dbConnect, t)

	err := CreateUser(c)

	if err == nil {
		t.Fatal("Expected error, but user was created")
	}
}
