package repository

import (
	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	repository "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

const createUserQuery string = `insert into "users"("username", "password") values($1, $2)`
const getUserByNameQuery string = "select id, username, password from users where username = $1"

var usersRepo usersRepository

// User repository
type usersRepository struct {
	repository.Worker
}

// Initialize users repository
func InitUserRepository(w *repository.Worker) {
	usersRepo = usersRepository{*w}
}

// Getter for users repository
func GetUsersRepo() *usersRepository {
	return &usersRepo
}

// Save new user
// Returns an error when the user cannot be saved, otherwise return nil
func (r *usersRepository) SaveUser(user model.User) error {
	_, err := r.Get().Exec(createUserQuery, user.UserName, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Get user by name
// Returns user and flag showing whether the user exists in the database
func (r *usersRepository) GetUserByUserName(userName string) (model.User, bool) {
	var user model.User
	err := r.Get().QueryRow(getUserByNameQuery, userName).Scan(&user.Id, &user.UserName, &user.Password)
	if err != nil {
		return model.User{}, false
	}
	return user, true
}
