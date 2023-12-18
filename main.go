package main

import (
	"bloggo/logging"
	"bloggo/templates"
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var conf = koanf.New(".")
var logger = logging.New()

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, highlighting.Highlighting),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

func toHtml(filename string) (string, error) {
	input, err := os.ReadFile(filename)

	if err != nil {
		return "", err
	}

	var html bytes.Buffer
	if err := md.Convert(input, &html); err != nil {
		logger.Errorf("Failed to convert %v", err)
	}

	return html.String(), nil
}

func dirToHtml(dir string) []string {
	files, err := os.ReadDir(dir)

	if err != nil {
		logger.Fatal(err)
	}

	var htmlFiles []string
	for _, file := range files {
		if !file.IsDir() {
			html, err := toHtml(filepath.Join(dir, file.Name()))

			if err != nil {
				logger.Fatal(err)
			}

			htmlFiles = append(htmlFiles, html)
		}
	}

	return htmlFiles
}

func main() {
	if err := conf.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		logger.Fatalf("Error loding config: %v", err)
	}

	server := echo.New()
	server.Use(logging.LogMiddleware(logger))

	server.GET("/", func(ctx echo.Context) error {
		return templates.Index().Render(context.Background(), ctx.Response().Writer)
	})

	posts := dirToHtml(conf.String("directories.posts"))

	server.GET("/post/:id", func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			log.Fatal("Failed to convert ID to string.")
		}

		return templates.Post(posts[id]).Render(context.Background(), ctx.Response().Writer)
	})

	logger.Fatal(server.Start(":1323"))
}
