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

type ArmyHandler struct {
	army *service.ArmyService
}

type trainUnitsRequest struct {
	UnitType string `json:"unitType"`
	Amount   int64  `json:"amount"`
}

type armyEnvelope struct {
	Army armyResponse `json:"army"`
}

type trainUnitsEnvelope struct {
	TrainingOrder trainingOrderResponse `json:"trainingOrder"`
	Resources     resourcesResponse     `json:"resources"`
}

type armyResponse struct {
	KingdomID      string                  `json:"kingdomId"`
	Units          []unitResponse          `json:"units"`
	TrainingOrders []trainingOrderResponse `json:"trainingOrders"`
	Summary        armySummaryResponse     `json:"summary"`
}

type unitResponse struct {
	Type           string                  `json:"type"`
	Label          string                  `json:"label"`
	Amount         int64                   `json:"amount"`
	Stats          unitStatsResponse       `json:"stats"`
	Cost           resourceValuesResponse  `json:"cost"`
	SecondsPerUnit int                     `json:"secondsPerUnit"`
	Requirements   unitRequirementResponse `json:"requirements"`
}

type unitStatsResponse struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
	Supply  int `json:"supply"`
}

type unitRequirementResponse struct {
	BarracksLevel int  `json:"barracksLevel"`
	IsMet         bool `json:"isMet"`
}

type trainingOrderResponse struct {
	ID          string     `json:"id"`
	UnitType    string     `json:"unitType"`
	UnitLabel   string     `json:"unitLabel"`
	Amount      int64      `json:"amount"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"startedAt"`
	FinishesAt  time.Time  `json:"finishesAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type armySummaryResponse struct {
	TotalUnits   int64 `json:"totalUnits"`
	TotalAttack  int64 `json:"totalAttack"`
	TotalDefense int64 `json:"totalDefense"`
	TotalSupply  int64 `json:"totalSupply"`
}

func NewArmyHandler(army *service.ArmyService) ArmyHandler {
	return ArmyHandler{army: army}
}

func (h ArmyHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	army, err := h.army.Current(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrArmyKingdomNotFound) {
			return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting army")
		}
		return err
	}

	return c.JSON(http.StatusOK, armyEnvelope{Army: newArmyResponse(army)})
}

func (h ArmyHandler) Train(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var request trainUnitsRequest
	if err := c.Bind(&request); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON")
	}

	result, err := h.army.Train(c.Request().Context(), userID, request.UnitType, request.Amount)
	if err != nil {
		return h.handleTrainError(c, err)
	}

	return c.JSON(http.StatusOK, trainUnitsEnvelope{
		TrainingOrder: newTrainingOrderResponse(result.Order),
		Resources:     newResourcesResponse(result.Resources),
	})
}

func (h ArmyHandler) handleTrainError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrArmyKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before training units")
	case errors.Is(err, service.ErrInvalidUnitType):
		return JSONError(c, http.StatusBadRequest, "invalid_unit_type", "Unit type is invalid")
	case errors.Is(err, service.ErrInvalidTrainingCount):
		return JSONError(c, http.StatusBadRequest, "invalid_training_amount", "Training amount must be between 1 and 50")
	case errors.Is(err, service.ErrInsufficientResources):
		return JSONError(c, http.StatusConflict, "insufficient_resources", "Not enough resources")
	case errors.Is(err, service.ErrBarracksLevelTooLow):
		return JSONError(c, http.StatusConflict, "barracks_level_too_low", "Barracks level is too low")
	default:
		return err
	}
}

func newArmyResponse(view service.ArmyView) armyResponse {
	units := make([]unitResponse, 0, len(view.Units))
	for _, unit := range view.Units {
		units = append(units, newUnitResponse(unit))
	}

	orders := make([]trainingOrderResponse, 0, len(view.TrainingOrders))
	for _, order := range view.TrainingOrders {
		orders = append(orders, newTrainingOrderResponse(order))
	}

	return armyResponse{
		KingdomID:      view.KingdomID,
		Units:          units,
		TrainingOrders: orders,
		Summary: armySummaryResponse{
			TotalUnits:   view.Summary.TotalUnits,
			TotalAttack:  view.Summary.TotalAttack,
			TotalDefense: view.Summary.TotalDefense,
			TotalSupply:  view.Summary.TotalSupply,
		},
	}
}

func newUnitResponse(view service.UnitView) unitResponse {
	return unitResponse{
		Type:           view.Unit.Type,
		Label:          view.Label,
		Amount:         view.Unit.Amount,
		Stats:          newUnitStatsResponse(view.Stats),
		Cost:           newResourceValuesResponse(view.Cost),
		SecondsPerUnit: view.Seconds,
		Requirements: unitRequirementResponse{
			BarracksLevel: view.Requirements.BarracksLevel,
			IsMet:         view.Requirements.IsMet,
		},
	}
}

func newUnitStatsResponse(stats gameconfig.UnitStats) unitStatsResponse {
	return unitStatsResponse{
		Attack:  stats.Attack,
		Defense: stats.Defense,
		Speed:   stats.Speed,
		Supply:  stats.Supply,
	}
}

func newTrainingOrderResponse(view service.TrainingOrderView) trainingOrderResponse {
	order := view.Order
	return trainingOrderResponse{
		ID:          order.ID,
		UnitType:    order.UnitType,
		UnitLabel:   view.UnitLabel,
		Amount:      order.Amount,
		Status:      order.Status,
		StartedAt:   order.StartedAt,
		FinishesAt:  order.FinishesAt,
		CompletedAt: order.CompletedAt,
	}
}
