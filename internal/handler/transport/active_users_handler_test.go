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
	var testDataUrl = map[string]string{"url": "http://localhost:8080/v1/user/active"}
	var testDataUrlM, _ = json.Marshal(testDataUrl)
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.POST, "/v1/chat/ws.rtm.start", bytes.NewReader(testDataUrlM))
	c := e.NewContext(req, rec)
	repository.InitActiveUsersStorage()

	err := GetActiveUsers(c)
	if err != nil {
		t.Fatalf("Failed to obtain active users")
	}
}
