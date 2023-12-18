package main

import (
	"bloggo/logging"

	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()
	logger := logging.New()

	server.Use(logging.LogMiddleware(logger))

	server.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, world!")
	})

	logger.Fatal(server.Start(":1323"))
}
