package handlers

import (
	"errors"
	"net/http"

	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type MeHandler struct {
	auth *service.AuthService
}

type meResponse struct {
	User userResponse `json:"user"`
}

func NewMeHandler(auth *service.AuthService) MeHandler {
	return MeHandler{auth: auth}
}

func (h MeHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	user, err := h.auth.CurrentUser(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return JSONError(c, http.StatusUnauthorized, "user_not_found", "User not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, meResponse{User: newUserResponse(user)})
}
