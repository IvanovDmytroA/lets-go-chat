package repository

import (
	"context"
	"sync"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	repository "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

var usersRepo usersRepository

// User repository
type usersRepository struct {
	W  repository.Worker
	mu sync.Mutex
}

// Initialize users repository
func InitUserRepository(w *repository.Worker) {
	usersRepo = usersRepository{W: *w}
}

// Getter for users repository
func GetUsersRepo() *usersRepository {
	return &usersRepo
}

// Save new user
// Returns an error when the user cannot be saved, otherwise return nil
func (ur *usersRepository) SaveUser(user model.User) error {
	ctx := context.Background()
	ur.mu.Lock()
	defer ur.mu.Unlock()
	_, err := ur.W.Get().NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Get user by name
// Returns user and flag showing whether the user exists in the database
func (ur *usersRepository) GetUserByUserName(userName string) (model.User, bool) {
	ctx := context.Background()
	var user model.User
	if err := ur.W.Get().NewSelect().Model(&user).Where("user_name = ?", userName).Scan(ctx); err != nil {
		return model.User{}, false
	}
	return user, true
}

// Delete user
// Returns an error when the user cannot be deleted, otherwise return nil
func (ur *usersRepository) DeleteUser(user model.User) error {
	ctx := context.Background()
	ur.mu.Lock()
	defer ur.mu.Unlock()
	_, err := ur.W.Get().NewDelete().Model(&user).Where("user_name = ?", &user.UserName).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
