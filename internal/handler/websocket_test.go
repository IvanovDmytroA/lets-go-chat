package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"

	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	td "github.com/IvanovDmytroA/lets-go-chat/tests"
)

var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjUxOGM3ZDc1LTA5MjAtNDY4ZC05MjhkLWViNTlmZTVjZWQ3NSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTYzOTY2MDY3NSwidXNlcl9pZCI6IjIwZDhmZDU2LWZiYTEtNDEwZS05ZTFmLTE1ZGEyNDE1NDdkNSJ9.uRSi_1grLTTRqg-odRUyMJTYilREAIU9i0eUZrGPm3M"
var userId string = "7efbd1e7-7e7c-4c3e-9928-85b17c5d9978"

func TestConnectInvalidConnectionType(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.CONNECT, "/v1/chat/ws.rtm.start?token="+token, bytes.NewReader([]byte{}))
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	keys := redisConnect.Keys("*").Val()
	for _, v := range keys {
		redisConnect.Del(v)
	}

	at := time.Unix(time.Now().Add(time.Minute*15).Unix(), 0)
	now := time.Now()
	errAccess := redisConnect.Set(token, userId, at.Sub(now)).Err()

	if errAccess != nil {
		t.Fatal("Failed to create redis auth")
	}

	err := connect(c, AccessDetails{AccessUuid: token, UserId: userId})

	keys = redisConnect.Keys("*").Val()
	for _, v := range keys {
		redisConnect.Del(v)
	}

	if err == nil {
		t.Fatal("Connected with invalid connection type")
	}

}

func TestConnectUnauthorized(t *testing.T) {
	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.CONNECT, "/v1/chat/ws.rtm.start?token="+token, bytes.NewReader([]byte{}))
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	keys := redisConnect.Keys("*").Val()
	for _, v := range keys {
		redisConnect.Del(v)
	}

	err := connect(c, AccessDetails{AccessUuid: token, UserId: userId})

	if err == nil {
		t.Fatal("Connected unauthorized user")
	}

}

func TestWebsocketInvalidToken(t *testing.T) {
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

	var testDataUrl = map[string]string{"url": "http://localhost:8080/v1/chat/ws.rtm.start?token=abcd"}
	var testDataUrlM, _ = json.Marshal(testDataUrl)

	req := httptest.NewRequest(echo.POST, "/v1/chat/ws.rtm.start", bytes.NewReader(testDataUrlM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	err = Websocket(c)

	keys := redisConnect.Keys("*").Val()
	for _, v := range keys {
		redisConnect.Del(v)
	}

	if err == nil {
		t.Fatal("Invalid token does not caused an error")
	}

}

func TestUpdateOnline(t *testing.T) {
	repository.InitActiveUsersStorage()

	updateOnline("user", true)

	if len(repository.GetActiveUsersStorage().Users) == 0 {
		t.Fatal("Active users list was not updated")
	}

	updateOnline("user", false)

	if len(repository.GetActiveUsersStorage().Users) != 0 {
		t.Fatal("Active users list was not updated")
	}
}
