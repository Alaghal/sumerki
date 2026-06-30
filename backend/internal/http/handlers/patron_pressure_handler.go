package handlers

import (
	"errors"
	"net/http"
	"time"

	"sumerki/backend/internal/gameconfig"
	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type PatronPressureHandler struct {
	pressure *service.PatronPressureService
}

type patronPressureEnvelope struct {
	Pressure *patronPressureResponse `json:"pressure"`
}

type patronPressurePaymentEnvelope struct {
	Pressure  patronPressureResponse `json:"pressure"`
	Resources resourcesResponse      `json:"resources"`
}

type patronPressureCrisisEnvelope struct {
	Patron   *patronRelationResponse `json:"patron,omitempty"`
	Pressure *patronPressureResponse `json:"pressure"`
	Kingdom  *patronKingdomResponse  `json:"kingdom,omitempty"`
}

type patronPressureRequest struct {
	Choice string `json:"choice"`
}

type patronPressureResponse struct {
	Patron            string                          `json:"patron"`
	PatronLabel       string                          `json:"patronLabel"`
	PressureLevel     int                             `json:"pressureLevel"`
	CrisisStatus      string                          `json:"crisisStatus"`
	TributeDebt       patronTributeDebtResponse       `json:"tributeDebt"`
	ContributionDebt  patronContributionDebtResponse  `json:"contributionDebt"`
	NextTributeAt     *time.Time                      `json:"nextTributeAt"`
	DelayUntil        *time.Time                      `json:"delayUntil"`
	Summary           string                          `json:"summary"`
	AvailableActions  []string                        `json:"availableActions"`
	ProtectedMinimums patronProtectedMinimumsResponse `json:"protectedMinimums"`
}

type patronTributeDebtResponse struct {
	Gold int64 `json:"gold"`
	Food int64 `json:"food"`
}

type patronContributionDebtResponse struct {
	Food int64 `json:"food"`
}

type patronProtectedMinimumsResponse struct {
	Gold       int64  `json:"gold"`
	Food       int64  `json:"food"`
	Wood       int64  `json:"wood"`
	Stone      int64  `json:"stone"`
	Population *int64 `json:"population"`
}

func NewPatronPressureHandler(pressure *service.PatronPressureService) PatronPressureHandler {
	return PatronPressureHandler{pressure: pressure}
}

func (h PatronPressureHandler) Current(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	result, err := h.pressure.Current(c.Request().Context(), userID)
	if err != nil {
		return patronPressureError(c, err)
	}

	var pressure *patronPressureResponse
	if result.Pressure != nil {
		response := newPatronPressureResponse(*result.Pressure)
		pressure = &response
	}
	return c.JSON(http.StatusOK, patronPressureEnvelope{Pressure: pressure})
}

func (h PatronPressureHandler) PayTribute(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	result, err := h.pressure.PayTribute(c.Request().Context(), userID)
	if err != nil {
		return patronPressureError(c, err)
	}

	return c.JSON(http.StatusOK, patronPressurePaymentEnvelope{
		Pressure:  newPatronPressureResponse(result.Pressure),
		Resources: newResourcesResponse(result.Resources),
	})
}

func (h PatronPressureHandler) ChooseCrisis(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var req patronPressureRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}

	result, err := h.pressure.ChooseCrisis(c.Request().Context(), userID, req.Choice)
	if err != nil {
		return patronPressureError(c, err)
	}

	var pressure *patronPressureResponse
	if result.Pressure != nil {
		response := newPatronPressureResponse(*result.Pressure)
		pressure = &response
	}
	var kingdom *patronKingdomResponse
	if result.Kingdom != nil {
		response := newPatronKingdomResponse(*result.Kingdom)
		kingdom = &response
	}
	return c.JSON(http.StatusOK, patronPressureCrisisEnvelope{
		Patron:   nil,
		Pressure: pressure,
		Kingdom:  kingdom,
	})
}

func newPatronPressureResponse(view service.PatronPressureView) patronPressureResponse {
	state := view.State
	nextTributeAt := &state.NextTributeAt
	if state.Patron == "independent" {
		nextTributeAt = nil
	}
	return patronPressureResponse{
		Patron:            state.Patron,
		PatronLabel:       view.PatronLabel,
		PressureLevel:     state.PressureLevel,
		CrisisStatus:      state.CrisisStatus,
		TributeDebt:       patronTributeDebtResponse{Gold: state.TributeDebtGold, Food: state.TributeDebtFood},
		ContributionDebt:  patronContributionDebtResponse{Food: state.ContributionDebtFood},
		NextTributeAt:     nextTributeAt,
		DelayUntil:        state.DelayUntil,
		Summary:           view.Summary,
		AvailableActions:  view.AvailableActions,
		ProtectedMinimums: newPatronProtectedMinimumsResponse(view.ProtectedMinimums),
	}
}

func newPatronProtectedMinimumsResponse(values gameconfig.ResourceValues) patronProtectedMinimumsResponse {
	return patronProtectedMinimumsResponse{
		Gold:       values.Gold,
		Food:       values.Food,
		Wood:       values.Wood,
		Stone:      values.Stone,
		Population: nil,
	}
}

func patronPressureError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrPatronKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting patron pressure")
	case errors.Is(err, service.ErrNoPatronSelected):
		return JSONError(c, http.StatusBadRequest, "no_patron_selected", "Choose a patron before using patron pressure")
	case errors.Is(err, service.ErrNoTributeDue):
		return JSONError(c, http.StatusBadRequest, "no_tribute_due", "No tribute or contribution is due")
	case errors.Is(err, service.ErrInvalidCrisisChoice):
		return JSONError(c, http.StatusBadRequest, "invalid_crisis_choice", "Invalid crisis choice")
	case errors.Is(err, service.ErrCrisisChoiceUnavailable):
		return JSONError(c, http.StatusBadRequest, "crisis_choice_not_available", "This crisis choice is not available")
	default:
		return err
	}
}
