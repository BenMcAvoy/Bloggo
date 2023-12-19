package main

import (
	"bloggo/logging"
	"bloggo/templates"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

var conf = koanf.New(".")
var logger = logging.New()

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, highlighting.NewHighlighting(
		highlighting.WithStyle("doom-one"),
		highlighting.WithFormatOptions(
			chromahtml.WithLineNumbers(true),
		),
	)),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

func toHtml(filename string) (map[string]string, error) {
	input, err := os.ReadFile(filename)

	if err != nil {
		logger.Fatal(err)
	}

	var html bytes.Buffer
	if err := md.Convert(input, &html); err != nil {
		logger.Errorf("Failed to convert %v", err)
	}

	return map[string]string{filename: html.String()}, nil
}

func dirToHtml(dir string) map[string]string {
	files, err := os.ReadDir(dir)

	if err != nil {
		logger.Fatal(err)
	}

	var htmlFiles = make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			html, err := toHtml(filepath.Join(dir, file.Name()))

			if err != nil {
				logger.Fatal(err)
			}

			for filename, content := range html {
				pat := filepath.Base(filename)
				ext := filepath.Ext(filename)

				key := strings.TrimSuffix(pat, ext)

				htmlFiles[key] = content
			}
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

	server.Static("/static", "static")

	server.GET("/", func(ctx echo.Context) error {
		return templates.Index().Render(context.Background(), ctx.Response().Writer)
	})

	posts := dirToHtml(conf.String("directories.posts"))

	server.GET("/post/:id", func(ctx echo.Context) error {
		id := ctx.Param("id")
		return templates.Base(posts[id]).Render(context.Background(), ctx.Response().Writer)
	})

	logger.Fatal(server.Start(":1323"))
}
