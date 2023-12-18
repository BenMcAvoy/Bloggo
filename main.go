package main

import (
	"bloggo/logging"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	hello "bloggo/templates"
)

func main() {
	server := echo.New()
	logger := logging.New()

	server.Use(logging.LogMiddleware(logger))

	server.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, world!")
	})

	server.GET("/hello/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")
		return hello.Hello(name).Render(context.Background(), ctx.Response().Writer)
	})

	logger.Fatal(server.Start(":1323"))
}
