package handlers

import (
	"sumerki/backend/internal/http/apierror"

	"github.com/labstack/echo/v4"
)

func JSONError(c echo.Context, status int, code string, message string) error {
	return apierror.JSON(c, status, code, message)
}
