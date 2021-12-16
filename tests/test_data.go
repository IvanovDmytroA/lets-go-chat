package test_data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/IvanovDmytroA/lets-go-chat/internal/model"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	connectors "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var TestData = map[string]string{"userName": "user", "password": "password"}
var IncompleteTestData = map[string]string{"pass": "pass"}
var TestDataM, _ = json.Marshal(TestData)
var IncompleteTestDataM, _ = json.Marshal(IncompleteTestData)
var userUuid = uuid.NewV4()
var user = model.User{
	Id:       userUuid.String(),
	UserName: "user",
	Password: "$2a$14$n887vQBDFPzCSmLbUzZZ4uwC9NKhE7guYLijT2Qw1Ne/SUFWqV4oO",
}

func DBConnection() *bun.DB {
	var worker connectors.Worker = &connectors.PostgresWorker{}
	var dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", "postgres", "pass", "localhost", 5432)
	}

	dbc := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db := bun.NewDB(dbc, pgdialect.New())
	worker.Init(db)
	repository.InitUserRepository(&worker)
	return db
}

func RedisConnection() *redis.Client {
	var redisUrl = os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = fmt.Sprintf("redis://%s:%d", "localhost", 6379)
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client

}

func AddTestUser(dbc *bun.DB) error {
	userRep := repository.GetUsersRepo()
	err := userRep.SaveUser(user)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func DelUser(dbConnect *bun.DB, t *testing.T) {
	ctx := context.Background()
	user := new(model.User)

	_, err := dbConnect.NewDelete().Model(user).Where("user_name = ?", TestData["userName"]).Exec(ctx)
	if err != nil {
		t.Fatalf("Error while deleting test user: %v", err)
	}
}
