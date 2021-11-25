package repository

import (
	"context"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	repository "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

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
	ctx := context.Background()
	_, err := r.Worker.Get().NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Get user by name
// Returns user and flag showing whether the user exists in the database
func (r *usersRepository) GetUserByUserName(userName string) (model.User, bool) {
	ctx := context.Background()
	var user model.User
	if err := r.Worker.Get().NewSelect().Model(&user).Where("user_name = ?", userName).Scan(ctx); err != nil {
		return model.User{}, false
	}
	return user, true
}
