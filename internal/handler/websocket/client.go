package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub  *hub
	conn *websocket.Conn
	send chan []byte
}

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

	return connect(c, *accessDetails)
}

func connect(c echo.Context, ad AccessDetails) error {
	redClient, _ := c.Get("redis").(*redis.Client)

	userId, err := redClient.Get(ad.AccessUuid).Result()
	if err != nil {
		errMsg := "Unauthorized. "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	updateOnline(userId, true)
	redClient.Del(ad.AccessUuid)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	hub := GetHub()
	client := &Client{hub: hub, conn: ws, send: make(chan []byte, 256)}
	client.hub.register <- client
	err = writePrevMessages(ws)
	if err != nil {
		log.Printf("Failed to download prev messages: %v", err)
	}

	go client.writePump(userId)
	go client.readPump(userId)

	return nil

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

func (c *Client) readPump(userId string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	wg := sync.WaitGroup{}
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				updateOnline(userId, false)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		wg.Add(1)
		c.hub.broadcast <- message
		go func(msg []byte, userId string) {
			repository.GetMessagesRepository().SaveMessage(model.Message{UserId: userId, Text: string(msg)})
			wg.Done()
		}(message, userId)
	}

	wg.Wait()
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(userId string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				updateOnline(userId, false)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				updateOnline(userId, false)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				updateOnline(userId, false)
				return
			}
		}
	}
}

func updateOnline(username string, online bool) {
	if online {
		repository.GetActiveUsersStorage().AddUserToActiveUsersList(username)
	} else {
		repository.GetActiveUsersStorage().RemoveUserFromActiveUsersList(username)
	}
}

func writePrevMessages(ws *websocket.Conn) error {
	messages, err := repository.GetMessagesRepository().GetAllMessages()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, msg := range messages {
		ws.WriteMessage(websocket.TextMessage, []byte(msg.Text))
	}

	return nil
}
