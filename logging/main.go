package logging

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func LogMiddleware(logger *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			err := next(c)

			if err != nil {
				logger.Error("Request", req.Method, req.URL.Path, err)
			} else {
				logger.Info("Request", req.Method, req.URL.Path, res.Status, time.Since(start))
			}

			return err
		}
	}
}

func New() *log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: "10",
			Dark:  "2",
		}).
		Foreground(lipgloss.Color("0"))

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERRO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: "9",
			Dark:  "1",
		}).
		Foreground(lipgloss.Color("0"))

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATA").
		Padding(0, 1, 0, 1).
		Background(lipgloss.AdaptiveColor{
			Light: "13",
			Dark:  "5",
		}).
		Foreground(lipgloss.Color("0"))
	logger := log.New(os.Stderr)
	logger.SetStyles(styles)

	return logger
}
