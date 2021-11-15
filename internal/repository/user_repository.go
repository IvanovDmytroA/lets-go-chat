package repository

import (
	"errors"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
)

var userStorage = make(map[string]model.User)

var ErrUserNotFound = errors.New("user doesn't exist")

func SaveUser(user model.User) {
	userStorage[user.UserName] = user
}

func GetUserByUserName(login string) (model.User, error) {
	if user, exists := userStorage[login]; exists {
		return user, nil
	}
	return model.User{}, ErrUserNotFound
}
