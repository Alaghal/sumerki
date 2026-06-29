package handlers

import (
	"errors"
	"net/http"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type PatronHandler struct {
	patrons *service.PatronService
}

type patronOptionResponse struct {
	Key              string   `json:"key"`
	Label            string   `json:"label"`
	ShortDescription string   `json:"shortDescription"`
	Flavor           string   `json:"flavor"`
	CurrentEffects   []string `json:"currentEffects"`
	FutureEffects    []string `json:"futureEffects"`
}

type patronOptionsEnvelope struct {
	Patrons []patronOptionResponse `json:"patrons"`
}

type patronStatusEnvelope struct {
	Patron           *patronRelationResponse `json:"patron"`
	AvailablePatrons []string                `json:"availablePatrons"`
}

type patronJoinRequest struct {
	Patron string `json:"patron"`
}

type patronJoinEnvelope struct {
	Patron  patronRelationResponse `json:"patron"`
	Kingdom patronKingdomResponse  `json:"kingdom"`
}

type patronBreakEnvelope struct {
	Patron  *patronRelationResponse `json:"patron"`
	Kingdom patronKingdomResponse   `json:"kingdom"`
}

type patronKingdomResponse struct {
	ID     string  `json:"id"`
	Patron *string `json:"patron"`
}

type patronRelationResponse struct {
	ID             string     `json:"id"`
	KingdomID      string     `json:"kingdomId"`
	Key            string     `json:"key"`
	Label          string     `json:"label"`
	Favor          int        `json:"favor"`
	Standing       string     `json:"standing"`
	JoinedAt       time.Time  `json:"joinedAt"`
	LeftAt         *time.Time `json:"leftAt"`
	CurrentEffects []string   `json:"currentEffects"`
	FutureEffects  []string   `json:"futureEffects"`
}

func NewPatronHandler(patrons *service.PatronService) PatronHandler {
	return PatronHandler{patrons: patrons}
}

func (h PatronHandler) Options(c echo.Context) error {
	if _, ok := appmiddleware.UserID(c); !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	options := h.patrons.Options()
	response := make([]patronOptionResponse, 0, len(options))
	for _, option := range options {
		response = append(response, newPatronOptionResponse(option))
	}
	return c.JSON(http.StatusOK, patronOptionsEnvelope{Patrons: response})
}

func (h PatronHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	status, err := h.patrons.Current(c.Request().Context(), userID)
	if err != nil {
		return patronError(c, err)
	}

	var patron *patronRelationResponse
	if status.Patron != nil {
		response := newPatronRelationResponse(*status.Patron)
		patron = &response
	}
	return c.JSON(http.StatusOK, patronStatusEnvelope{
		Patron:           patron,
		AvailablePatrons: status.AvailablePatrons,
	})
}

func (h PatronHandler) Join(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var req patronJoinRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	result, err := h.patrons.Join(c.Request().Context(), userID, req.Patron)
	if err != nil {
		return patronError(c, err)
	}

	return c.JSON(http.StatusOK, patronJoinEnvelope{
		Patron:  newPatronRelationResponse(result.Patron),
		Kingdom: newPatronKingdomResponse(result.Kingdom),
	})
}

func (h PatronHandler) Break(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	result, err := h.patrons.Break(c.Request().Context(), userID)
	if err != nil {
		return patronError(c, err)
	}

	return c.JSON(http.StatusOK, patronBreakEnvelope{
		Patron:  nil,
		Kingdom: newPatronKingdomResponse(result.Kingdom),
	})
}

func newPatronOptionResponse(option gameconfig.PatronConfig) patronOptionResponse {
	return patronOptionResponse{
		Key:              option.Key,
		Label:            option.Label,
		ShortDescription: option.ShortDescription,
		Flavor:           option.Flavor,
		CurrentEffects:   option.CurrentEffects,
		FutureEffects:    option.FutureEffects,
	}
}

func newPatronRelationResponse(view service.PatronRelationView) patronRelationResponse {
	relation := view.Relation
	return patronRelationResponse{
		ID:             relation.ID,
		KingdomID:      relation.KingdomID,
		Key:            relation.Patron,
		Label:          view.Label,
		Favor:          relation.Favor,
		Standing:       relation.Standing,
		JoinedAt:       relation.JoinedAt,
		LeftAt:         relation.LeftAt,
		CurrentEffects: view.CurrentEffects,
		FutureEffects:  view.FutureEffects,
	}
}

func newPatronKingdomResponse(kingdom domain.Kingdom) patronKingdomResponse {
	return patronKingdomResponse{
		ID:     kingdom.ID,
		Patron: kingdom.Patron,
	}
}

func patronError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrPatronKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting patron status")
	case errors.Is(err, service.ErrInvalidPatron):
		return JSONError(c, http.StatusBadRequest, "invalid_patron", "Invalid patron")
	default:
		return err
	}
}
