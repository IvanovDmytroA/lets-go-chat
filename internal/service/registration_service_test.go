package service

import (
	"bytes"
	"net/http"
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

	ur, status := CreateUser("user", "password")

	if len(ur.Id) == 0 {
		t.Fatalf("Created user does not containt generated ID")
	}
	if len(ur.UserName) == 0 {
		t.Fatalf("Created user does not containt username")
	}
	if status != http.StatusOK {
		t.Fatalf("Status is not 200 OK")
	}

	td.DelUser(dbConnect, t)
}

func TestCreateUserWithEmptyCredentials(t *testing.T) {
	dbConnect := td.DBConnection()
	td.DelUser(dbConnect, t)

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/user", bytes.NewReader(td.TestDataM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)

	ur, status := CreateUser("", "")

	if len(ur.Id) != 0 || status == http.StatusOK {
		t.Fatalf("Created user with empty credentials")
	}

	td.DelUser(dbConnect, t)
}
