package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

func Websocket(c echo.Context) error {
	userToken := c.QueryParam("token")
	accessDetails, err := extractTokenMetadata(userToken)
	if err != nil {
		errMsg := "Invalid token: "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	client, _ := c.Get("redis").(*redis.Client)

	userId, err := client.Get(accessDetails.AccessUuid).Result()
	if err != nil {
		errMsg := "Unauthorized. "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	updateOnline(userId, true)
	client.Del(accessDetails.AccessUuid)

	upgrader := websocket.Upgrader{}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Read message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			updateOnline(userId, false)
			return nil // todo: implement closed connection error hsndling
		}

		// Write echo message
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			updateOnline(userId, false)
			return nil // todo: implement closed connection error hsndling
		}
	}
}

func verifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractTokenMetadata(t string) (*AccessDetails, error) {
	vt, err := verifyToken(t)
	if err != nil {
		return nil, err
	}
	claims, ok := vt.Claims.(jwt.MapClaims)
	if ok && vt.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := claims["user_id"].(string)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func updateOnline(username string, online bool) {
	if online {
		repository.GetActiveUsersStorage().AddUserToActiveUsersList(username)
	} else {
		repository.GetActiveUsersStorage().RemoveUserFromActiveUsersList(username)
	}
}
