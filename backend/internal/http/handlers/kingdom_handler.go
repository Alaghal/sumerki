package handlers

import (
	"errors"
	"net/http"
	"time"

	"sumerki/backend/internal/domain"
	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type KingdomHandler struct {
	kingdoms *service.KingdomService
}

type createKingdomRequest struct {
	Name    string `json:"name"`
	Culture string `json:"culture"`
}

type kingdomResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Culture   string    `json:"culture"`
	Patron    *string   `json:"patron"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type kingdomEnvelope struct {
	Kingdom *kingdomResponse `json:"kingdom"`
}

func NewKingdomHandler(kingdoms *service.KingdomService) KingdomHandler {
	return KingdomHandler{kingdoms: kingdoms}
}

func (h KingdomHandler) Create(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var req createKingdomRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	kingdom, err := h.kingdoms.Create(c.Request().Context(), userID, req.Name, req.Culture)
	if err != nil {
		return h.handleCreateError(c, err)
	}

	response := newKingdomResponse(kingdom)
	return c.JSON(http.StatusCreated, kingdomEnvelope{Kingdom: &response})
}

func (h KingdomHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	kingdom, err := h.kingdoms.Current(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	if kingdom == nil {
		return c.JSON(http.StatusOK, kingdomEnvelope{Kingdom: nil})
	}

	response := newKingdomResponse(*kingdom)
	return c.JSON(http.StatusOK, kingdomEnvelope{Kingdom: &response})
}

func (h KingdomHandler) handleCreateError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrKingdomNameTooShort):
		return JSONError(c, http.StatusBadRequest, "kingdom_name_too_short", "Kingdom name must be at least 3 characters")
	case errors.Is(err, service.ErrKingdomNameTooLong):
		return JSONError(c, http.StatusBadRequest, "kingdom_name_too_long", "Kingdom name must be at most 32 characters")
	case errors.Is(err, service.ErrInvalidCulture):
		return JSONError(c, http.StatusBadRequest, "invalid_culture", "Culture is invalid")
	case errors.Is(err, service.ErrKingdomAlreadyExists):
		return JSONError(c, http.StatusConflict, "kingdom_already_exists", "Kingdom already exists")
	default:
		return err
	}
}

func newKingdomResponse(kingdom domain.Kingdom) kingdomResponse {
	return kingdomResponse{
		ID:        kingdom.ID,
		UserID:    kingdom.UserID,
		Name:      kingdom.Name,
		Culture:   kingdom.Culture,
		Patron:    kingdom.Patron,
		CreatedAt: kingdom.CreatedAt,
		UpdatedAt: kingdom.UpdatedAt,
	}
}
