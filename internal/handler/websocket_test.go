package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	td "github.com/IvanovDmytroA/lets-go-chat/tests"
)

func TestWebsocket(t *testing.T) {
	var TestDataUrl = map[string]string{"url": "http://localhost:8080/v1/chat/ws.rtm.start?token=7efbd1e7-7e7c-4c3e-9928-85b17c5d9978"}
	var TestDataUrlM, _ = json.Marshal(TestDataUrl)

	dbConnect := td.DBConnection()
	redisConnect := td.RedisConnection()

	td.DelUser(dbConnect, t)
	err := td.AddTestUser(dbConnect)
	if err != nil {
		t.Fatalf("Saving test user failed")
	}

	e := echo.New()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/chat/ws.rtm.start", bytes.NewReader(TestDataUrlM))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := e.NewContext(req, rec)

	c.Set("db", dbConnect)
	c.Set("redis", redisConnect)

	errWS := Websocket(c)
	invalidTokenErrMsg := "code=400, message=Invalid token: token contains an invalid number of segments"

	td.DelUser(dbConnect, t)

	if errWS != nil && errWS.Error() != invalidTokenErrMsg {
		t.Fatalf("Invalid token passed")
	}
}
