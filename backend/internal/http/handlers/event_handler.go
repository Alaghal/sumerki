package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	events *service.EventService
}

type eventsEnvelope struct {
	Events      []kingdomEventResponse `json:"events"`
	ActiveCount int                    `json:"activeCount"`
}

type chooseEventRequest struct {
	ChoiceKey string `json:"choiceKey"`
}

type chooseEventEnvelope struct {
	Event     kingdomEventResponse `json:"event"`
	Resources *resourcesResponse   `json:"resources,omitempty"`
	Army      *armyResponse        `json:"army,omitempty"`
	Kingdom   kingdomResponse      `json:"kingdom"`
}

type kingdomEventResponse struct {
	ID                string                `json:"id"`
	EventKey          string                `json:"eventKey"`
	Category          string                `json:"category"`
	Title             string                `json:"title"`
	Body              string                `json:"body"`
	Status            string                `json:"status"`
	GeneratedAt       time.Time             `json:"generatedAt"`
	ExpiresAt         time.Time             `json:"expiresAt"`
	ResolvedAt        *time.Time            `json:"resolvedAt"`
	SelectedChoiceKey *string               `json:"selectedChoiceKey"`
	Choices           []eventChoiceResponse `json:"choices"`
	Result            *service.EventResult  `json:"result"`
}

type eventChoiceResponse struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

func NewEventHandler(events *service.EventService) EventHandler {
	return EventHandler{events: events}
}

func (h EventHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	includeResolved := parseBool(c.QueryParam("includeResolved"), true)
	limit := parsePositiveInt(c.QueryParam("limit"), 20)
	if limit > 50 {
		limit = 50
	}

	result, err := h.events.Current(c.Request().Context(), userID, includeResolved, limit)
	if err != nil {
		return eventError(c, err)
	}

	events := make([]kingdomEventResponse, 0, len(result.Events))
	for _, event := range result.Events {
		events = append(events, newKingdomEventResponse(event))
	}
	return c.JSON(http.StatusOK, eventsEnvelope{
		Events:      events,
		ActiveCount: result.ActiveCount,
	})
}

func (h EventHandler) Choose(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var req chooseEventRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	result, err := h.events.Choose(c.Request().Context(), userID, c.Param("id"), req.ChoiceKey)
	if err != nil {
		return eventError(c, err)
	}

	var resources *resourcesResponse
	if result.Resources != nil {
		response := newResourcesResponse(*result.Resources)
		resources = &response
	}
	var army *armyResponse
	if result.Army != nil {
		response := newArmyResponse(*result.Army)
		army = &response
	}
	return c.JSON(http.StatusOK, chooseEventEnvelope{
		Event:     newKingdomEventResponse(result.Event),
		Resources: resources,
		Army:      army,
		Kingdom:   newKingdomResponse(result.Kingdom),
	})
}

func newKingdomEventResponse(view service.KingdomEventView) kingdomEventResponse {
	event := view.Event
	choices := make([]eventChoiceResponse, 0, len(view.Choices))
	for _, choice := range view.Choices {
		choices = append(choices, eventChoiceResponse{
			Key:         choice.Choice.Key,
			Label:       choice.Choice.Label,
			Description: choice.Choice.Description,
		})
	}
	return kingdomEventResponse{
		ID:                event.ID,
		EventKey:          event.GameEvent.Key,
		Category:          event.GameEvent.Category,
		Title:             event.GameEvent.Title,
		Body:              event.GameEvent.Body,
		Status:            event.Status,
		GeneratedAt:       event.GeneratedAt,
		ExpiresAt:         event.ExpiresAt,
		ResolvedAt:        event.ResolvedAt,
		SelectedChoiceKey: event.SelectedChoiceKey,
		Choices:           choices,
		Result:            view.Result,
	}
}

func parseBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func eventError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrEventKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting events")
	case errors.Is(err, service.ErrEventNotFound):
		return JSONError(c, http.StatusNotFound, "event_not_found", "Event not found")
	case errors.Is(err, service.ErrEventExpired):
		return JSONError(c, http.StatusConflict, "event_expired", "This event has expired")
	case errors.Is(err, service.ErrEventAlreadyResolved):
		return JSONError(c, http.StatusConflict, "event_already_resolved", "This event is already resolved")
	case errors.Is(err, service.ErrInvalidEventChoice):
		return JSONError(c, http.StatusBadRequest, "invalid_event_choice", "Invalid event choice")
	case errors.Is(err, service.ErrEventChoiceNotAvailable):
		return JSONError(c, http.StatusBadRequest, "event_choice_not_available", "Event choice is not available")
	default:
		return err
	}
}
