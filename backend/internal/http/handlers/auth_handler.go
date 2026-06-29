package handlers

import (
	"errors"
	"net/http"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	auth *service.AuthService
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewAuthHandler(auth *service.AuthService) AuthHandler {
	return AuthHandler{auth: auth}
}

func (h AuthHandler) Register(c echo.Context) error {
	var req authRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	result, err := h.auth.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return h.handleRegisterError(c, err)
	}

	return c.JSON(http.StatusCreated, newAuthResponse(result))
}

func (h AuthHandler) Login(c echo.Context) error {
	var req authRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	result, err := h.auth.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return JSONError(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password")
		}
		return err
	}

	return c.JSON(http.StatusOK, newAuthResponse(result))
}

func (h AuthHandler) handleRegisterError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidEmail):
		return JSONError(c, http.StatusBadRequest, "invalid_email", "Email is invalid")
	case errors.Is(err, service.ErrPasswordTooShort):
		return JSONError(c, http.StatusBadRequest, "password_too_short", "Password must be at least 8 characters")
	case errors.Is(err, service.ErrEmailAlreadyExists):
		return JSONError(c, http.StatusConflict, "email_already_exists", "Email already exists")
	default:
		return err
	}
}

func newAuthResponse(result service.AuthResult) authResponse {
	return authResponse{
		User:  newUserResponse(result.User),
		Token: result.Token,
	}
}

func newUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
