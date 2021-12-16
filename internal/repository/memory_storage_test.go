package repository

import (
	"testing"
)

func TestInitActiveUsersStorage(t *testing.T) {
	InitActiveUsersStorage()

	if GetActiveUsersStorage() == nil {
		t.Fatal("Failed to initialize memory users storage")
	}
}

func TestAddUserToActiveUsersList(t *testing.T) {
	InitActiveUsersStorage()
	storage := GetActiveUsersStorage()
	storage.AddUserToActiveUsersList("user")

	if len(storage.Users) == 0 {
		t.Fatal("Failed add user to active users list")
	}
}

func TestAddUserToActiveUsersListOnce(t *testing.T) {
	InitActiveUsersStorage()
	storage := GetActiveUsersStorage()
	storage.AddUserToActiveUsersList("user")
	storage.AddUserToActiveUsersList("user")

	if len(storage.Users) > 1 {
		t.Fatal("Failed add user to active users list")
	}
}

func TestRemoveUserFromActiveUsersList(t *testing.T) {
	InitActiveUsersStorage()
	storage := GetActiveUsersStorage()
	storage.AddUserToActiveUsersList("user")

	if len(storage.Users) == 0 {
		t.Fatal("Failed add user to active users list")
	}

	storage.RemoveUserFromActiveUsersList("user")

	if len(storage.Users) > 0 {
		t.Fatal("Failed delete user from active users list")
	}
}
