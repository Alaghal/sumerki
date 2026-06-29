package httpserver

import (
	"database/sql"
	"errors"
	"net/http"

	"sumerki/backend/internal/http/handlers"
	appmiddleware "sumerki/backend/internal/http/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(database *sql.DB) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = errorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(appmiddleware.LocalCORS()))

	health := handlers.NewHealth(database)
	e.GET("/health", health.Health)
	e.GET("/ready", health.Ready)

	return e
}

func errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError
	code := "internal_error"
	message := "Internal server error"

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		status = httpErr.Code
		message = http.StatusText(status)
		if message == "" {
			message = "HTTP error"
		}
		code = "http_error"
	}

	if writeErr := handlers.JSONError(c, status, code, message); writeErr != nil {
		c.Logger().Error(writeErr)
	}
}
