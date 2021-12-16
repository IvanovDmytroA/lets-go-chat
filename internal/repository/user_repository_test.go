package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	connectors "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
)

const userName string = "user"

var userModel model.User = model.User{Id: "ID", UserName: userName, Password: "password"}

func TestGetUserRepo(t *testing.T) {
	initDb()

	repo := GetUsersRepo()

	if repo == nil {
		t.Fatal("Failed to obtain uses repo")
	}
}

func TestSaveUser(t *testing.T) {
	initDb()

	repo := GetUsersRepo()
	err := repo.SaveUser(userModel)

	if err != nil {
		t.Fatalf("Failed to save user: %s", err.Error())
	}
}

func TestSaveUserReturnError(t *testing.T) {
	initDb()

	repo := GetUsersRepo()
	err := repo.SaveUser(model.User{})

	if err == nil {
		t.Fatal("Saved user without required fields")
	}
}

func TestGetUserByUserName(t *testing.T) {
	initDb()

	repo := GetUsersRepo()
	err := repo.SaveUser(userModel)

	if err != nil {
		t.Fatalf("Failed to save user: %s", err.Error())
	}

	user, exists := repo.GetUserByUserName(userName)

	if !exists {
		t.Fatalf("User with name %s does not exist", userName)
	}
	if user.UserName != userName {
		t.Fatalf("Returned wrong user: %s", user.UserName)
	}
}

func TestGetUserByUserNameReturnEmptyUser(t *testing.T) {
	initDb()

	repo := GetUsersRepo()
	_, exists := repo.GetUserByUserName(userName)

	if exists {
		t.Fatal("Returned not existent user")
	}
}

func initDb() {
	var worker connectors.Worker = &connectors.PostgresWorker{}
	var dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", "postgres", "pass", "localhost", 5432)
	}

	dbc := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db := bun.NewDB(dbc, pgdialect.New())
	worker.Init(db)
	InitUserRepository(&worker)
	repo := GetUsersRepo()
	repo.DeleteUser(userModel)
}
