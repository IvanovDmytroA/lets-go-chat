package server

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

func dataSourceMiddleware(dataStore *bun.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dataStore)
			return next(c)
		}
	}
}

func redisMiddleware(cl *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("redis", cl)
			return next(c)
		}
	}
}

func bodyDumpMiddleware(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("\nRequest Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}
