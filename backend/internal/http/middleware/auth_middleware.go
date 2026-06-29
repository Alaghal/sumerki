package middleware

import (
	"errors"
	"net/http"
	"strings"

	"sumerki/backend/internal/http/apierror"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

const userIDContextKey = "user_id"

func Auth(auth *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return apierror.JSON(c, http.StatusUnauthorized, "missing_authorization_header", "Authorization header is required")
			}

			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apierror.JSON(c, http.StatusUnauthorized, "invalid_authorization_header", "Authorization header must be Bearer token")
			}

			userID, err := auth.ValidateToken(parts[1])
			if err != nil {
				if errors.Is(err, service.ErrExpiredToken) {
					return apierror.JSON(c, http.StatusUnauthorized, "expired_token", "Token has expired")
				}
				return apierror.JSON(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
			}

			c.Set(userIDContextKey, userID)
			return next(c)
		}
	}
}

func UserID(c echo.Context) (string, bool) {
	userID, ok := c.Get(userIDContextKey).(string)
	return userID, ok && userID != ""
}
