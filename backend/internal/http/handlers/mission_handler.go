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

type MissionHandler struct {
	missions *service.MissionService
}

type startMissionRequest struct {
	MissionKey string                    `json:"missionKey"`
	Units      []startMissionUnitRequest `json:"units"`
}

type startMissionUnitRequest struct {
	UnitType string `json:"unitType"`
	Amount   int64  `json:"amount"`
}

type availableMissionsEnvelope struct {
	Missions []availableMissionResponse `json:"missions"`
}

type missionsEnvelope struct {
	Missions []missionResponse `json:"missions"`
}

type startMissionEnvelope struct {
	Mission missionResponse `json:"mission"`
	Army    armyResponse    `json:"army"`
}

type availableMissionResponse struct {
	Key                 string                      `json:"key"`
	Label               string                      `json:"label"`
	Type                string                      `json:"type"`
	Description         string                      `json:"description"`
	DurationSeconds     int                         `json:"durationSeconds"`
	MinimumRequirements missionRequirementsResponse `json:"minimumRequirements"`
	BaseRewards         resourceValuesResponse      `json:"baseRewards"`
	Risk                string                      `json:"risk"`
}

type missionRequirementsResponse struct {
	TotalUnits int64 `json:"totalUnits"`
	Scouts     int64 `json:"scouts"`
}

type missionResponse struct {
	ID           string                 `json:"id"`
	MissionKey   string                 `json:"missionKey"`
	MissionLabel string                 `json:"missionLabel"`
	MissionType  string                 `json:"missionType"`
	Status       string                 `json:"status"`
	StartedAt    time.Time              `json:"startedAt"`
	FinishesAt   time.Time              `json:"finishesAt"`
	CompletedAt  *time.Time             `json:"completedAt"`
	Units        []missionUnitResponse  `json:"units"`
	Result       *missionResultResponse `json:"result"`
}

type missionUnitResponse struct {
	UnitType       string `json:"unitType"`
	UnitLabel      string `json:"unitLabel"`
	AmountSent     int64  `json:"amountSent"`
	AmountLost     int64  `json:"amountLost"`
	AmountReturned int64  `json:"amountReturned"`
}

type missionResultResponse struct {
	Result  string                 `json:"result"`
	Rewards resourceValuesResponse `json:"rewards"`
	Losses  map[string]int64       `json:"losses"`
}

func NewMissionHandler(missions *service.MissionService) MissionHandler {
	return MissionHandler{missions: missions}
}

func (h MissionHandler) Available(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	missions, err := h.missions.Available(c.Request().Context(), userID)
	if err != nil {
		return h.handleMissionError(c, err)
	}

	response := make([]availableMissionResponse, 0, len(missions))
	for _, mission := range missions {
		response = append(response, newAvailableMissionResponse(mission))
	}

	return c.JSON(http.StatusOK, availableMissionsEnvelope{Missions: response})
}

func (h MissionHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	missions, err := h.missions.Current(c.Request().Context(), userID)
	if err != nil {
		return h.handleMissionError(c, err)
	}

	response := make([]missionResponse, 0, len(missions))
	for _, mission := range missions {
		response = append(response, newMissionResponse(mission))
	}

	return c.JSON(http.StatusOK, missionsEnvelope{Missions: response})
}

func (h MissionHandler) Start(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var request startMissionRequest
	if err := c.Bind(&request); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON")
	}

	units := make([]service.StartMissionUnit, 0, len(request.Units))
	for _, unit := range request.Units {
		units = append(units, service.StartMissionUnit{
			UnitType: unit.UnitType,
			Amount:   unit.Amount,
		})
	}

	result, err := h.missions.Start(c.Request().Context(), userID, request.MissionKey, units)
	if err != nil {
		return h.handleMissionError(c, err)
	}

	return c.JSON(http.StatusOK, startMissionEnvelope{
		Mission: newMissionResponse(result.Mission),
		Army:    newArmyResponse(result.Army),
	})
}

func (h MissionHandler) handleMissionError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrMissionKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting missions")
	case errors.Is(err, service.ErrInvalidMissionKey):
		return JSONError(c, http.StatusBadRequest, "invalid_mission_key", "Mission key is invalid")
	case errors.Is(err, service.ErrInvalidMissionUnitType):
		return JSONError(c, http.StatusBadRequest, "invalid_unit_type", "Unit type is invalid")
	case errors.Is(err, service.ErrInvalidMissionUnitAmount):
		return JSONError(c, http.StatusBadRequest, "invalid_unit_amount", "Unit amount is invalid")
	case errors.Is(err, service.ErrInsufficientUnits):
		return JSONError(c, http.StatusConflict, "insufficient_units", "Not enough available units")
	case errors.Is(err, service.ErrMissionRequirementsNotMet):
		return JSONError(c, http.StatusConflict, "mission_requirements_not_met", "Mission requirements are not met")
	default:
		return err
	}
}

func newAvailableMissionResponse(mission gameconfig.MissionConfig) availableMissionResponse {
	return availableMissionResponse{
		Key:             mission.Key,
		Label:           mission.Label,
		Type:            mission.Type,
		Description:     mission.Description,
		DurationSeconds: mission.DurationSeconds,
		MinimumRequirements: missionRequirementsResponse{
			TotalUnits: mission.MinimumRequirements.TotalUnits,
			Scouts:     mission.MinimumRequirements.Scouts,
		},
		BaseRewards: newResourceValuesResponse(mission.BaseRewards),
		Risk:        mission.Risk,
	}
}

func newMissionResponse(view service.MissionView) missionResponse {
	mission := view.Mission
	units := make([]missionUnitResponse, 0, len(view.Units))
	for _, unit := range view.Units {
		units = append(units, missionUnitResponse{
			UnitType:       unit.Unit.UnitType,
			UnitLabel:      unit.Label,
			AmountSent:     unit.Unit.AmountSent,
			AmountLost:     unit.Unit.AmountLost,
			AmountReturned: unit.Unit.AmountReturned,
		})
	}

	return missionResponse{
		ID:           mission.ID,
		MissionKey:   mission.Key,
		MissionLabel: view.Label,
		MissionType:  mission.Type,
		Status:       mission.Status,
		StartedAt:    mission.StartedAt,
		FinishesAt:   mission.FinishesAt,
		CompletedAt:  mission.CompletedAt,
		Units:        units,
		Result:       newMissionResultResponse(view.Result),
	}
}

func newMissionResultResponse(result *service.MissionResult) *missionResultResponse {
	if result == nil {
		return nil
	}
	return &missionResultResponse{
		Result:  result.Result,
		Rewards: newResourceValuesResponse(result.Rewards),
		Losses:  result.Losses,
	}
}
