package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/labstack/echo"
)

func TestGetActiveUsers(t *testing.T) {
	var TestDataUrl = map[string]string{"url": "http://localhost:8080/v1/chat/ws.rtm.start?token=7efbd1e7-7e7c-4c3e-9928-85b17c5d9978"}
	var TestDataUrlM, _ = json.Marshal(TestDataUrl)
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/chat/ws.rtm.start", bytes.NewReader(TestDataUrlM))
	c := e.NewContext(req, rec)
	repository.InitActiveUsersStorage()

	err := GetActiveUsers(c)
	if err != nil {
		t.Fatalf("Failed to obtain active users")
	}
}
