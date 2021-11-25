package repository

import "log"

var activeUsers activeUsersStorage

type activeUsersStorage struct {
	Users []string
}

func InitActiveUsersStorage() {
	activeUsers = activeUsersStorage{}
}

func GetActiveUsersStorage() *activeUsersStorage {
	return &activeUsers
}

func (au *activeUsersStorage) AddUserToActiveUsersList(username string) {
	users := au.Users
	for _, ele := range users {
		if ele == username {
			log.Println("User already active")
			return
		}
	}
	au.Users = append(users, username)
}

func (au *activeUsersStorage) RemoveUserFromActiveUsersList(username string) {
	users := au.Users
	for index, ele := range users {
		if ele == username {
			au.Users = append(users[:index], users[index+1:]...)
			return
		}
	}
}
