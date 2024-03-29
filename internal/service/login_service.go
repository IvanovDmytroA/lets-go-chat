package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	"github.com/IvanovDmytroA/lets-go-chat/pkg/hasher"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

const errorMessage string = "Invalid name or password"

func LoginUser(userName, password string, c echo.Context) error {
	loginRequest := handler.LoginUserRequest{UserName: userName, Password: password}
	user, err := getUserFromRepo(loginRequest)
	if err != nil {
		fmt.Println(err)
		return err
	}
	token, err := createToken(user.Id)
	if err != nil {
		errMsg := "Failed to generate token"
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg+err.Error())
	}

	client, _ := c.Get("redis").(*redis.Client)
	saveErr := createAuth(client, user.Id, token)
	if saveErr != nil {
		errMsg := "CreateAuth failed: "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errMsg+err.Error())
	}

	loginUrl := "wss://serene-everglades-55494.herokuapp.com/v1/chat/ws.rtm.start?token=" + token.AccessToken
	loginUserResponse := handler.LoginUserResponse{
		Url: loginUrl,
	}
	enc := json.NewEncoder(c.Response())
	enc.SetEscapeHTML(false)
	err = enc.Encode(loginUserResponse)
	if err != nil {
		errMsg := "Failed to encode UserStorage: "
		log.Printf(errMsg+"%v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errMsg+err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().Header().Set("X-Rate-Limit", strconv.Itoa(360))
	c.Response().Header().Set("X-Expires-After", time.Now().Add(time.Minute*10).Format(time.RFC1123))
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	return nil
}

func getUserFromRepo(loginRequest handler.LoginUserRequest) (model.User, error) {
	userRepo := repository.GetUsersRepo()
	user, exists := userRepo.GetUserByUserName(loginRequest.UserName)
	if !exists {
		return user, echo.NewHTTPError(http.StatusBadRequest, "User with provided name does not exist")
	}

	isCorrectPassword := hasher.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isCorrectPassword {
		return user, echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	return user, nil
}

func createToken(userid string) (*handler.TokenDetails, error) {
	td := &handler.TokenDetails{}
	td.Expires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	var err error
	os.Setenv("ACCESS_SECRET", "")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.Expires
	atClaims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func createAuth(client *redis.Client, userid string, td *handler.TokenDetails) error {
	at := time.Unix(td.Expires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
