package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/IvanovDmytroA/lets-go-chat/internal/configuration"
	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	transport_handler "github.com/IvanovDmytroA/lets-go-chat/internal/handler/transport"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	connectors "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Start and configure server
func Start() {
	port := initServer()
	env, err := configuration.InitEnv()
	if err != nil {
		log.Fatal("Failed to init environment configuration")
	}

	initDb(env)
	repository.InitActiveUsersStorage()
	redisClient := initRedis(env)
	initEcho(redisClient, port)
}

func initServer() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func initDb(e *configuration.Env) {
	switch e.DataBase.Type {
	case "postgres":
		var worker connectors.Worker = &connectors.PostgresWorker{}
		var dbUrl = os.Getenv("DATABASE_URL")
		if dbUrl == "" {
			dbUrl = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
				e.DataBase.Host, e.DataBase.Port, e.DataBase.User, e.DataBase.Password)
		}

		dbc := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
		db := bun.NewDB(dbc, pgdialect.New())
		worker.Init(db)
		repository.InitUserRepository(&worker)
	}
}

func initRedis(e *configuration.Env) *redis.Client {
	var redisUrl = os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = fmt.Sprintf("redis://%s:%d", e.Redis.Host, e.Redis.Port)
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

func initEcho(rc *redis.Client, p string) {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDump(bodyDumpHandler))
	e.Use(middleware.Recover())
	s := handler.NewStatistic()
	e.Use(s.Process)
	e.Use(dataSourceMiddleware(repository.GetUsersRepo().Get()))
	e.Use(redisMiddleware(rc))
	e.GET("/stats", s.Handle)
	e.POST("/v1/user", transport_handler.CreateUser)
	e.POST("/v1/user/login", transport_handler.LoginUser)
	e.GET("/v1/user/active", transport_handler.GetActiveUsers)
	e.GET("/v1/chat/ws.rtm.start", handler.Websocket)
	e.Static("/v1/chat", "internal/public")

	e.Logger.Fatal(e.Start(":" + p))
}

// Middleware function to transfer sql dataStore to handlers
func dataSourceMiddleware(dataStore *bun.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dataStore)
			return next(c)
		}
	}
}

// Middleware function to transfer redis to handlers
func redisMiddleware(dataStore *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("redis", dataStore)
			return next(c)
		}
	}
}

// Middleware to log requests and responses
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("\nRequest Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}
