package handlers

import (
	"errors"
	"net/http"
	"time"

	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type BuildingHandler struct {
	buildings *service.BuildingService
}

type buildingResponse struct {
	ID                string               `json:"id"`
	KingdomID         string               `json:"kingdomId"`
	Type              string               `json:"type"`
	Label             string               `json:"label"`
	Level             int                  `json:"level"`
	MaxLevel          int                  `json:"maxLevel"`
	IsUpgrading       bool                 `json:"isUpgrading"`
	UpgradeStartedAt  *time.Time           `json:"upgradeStartedAt"`
	UpgradeFinishesAt *time.Time           `json:"upgradeFinishesAt"`
	NextUpgrade       *nextUpgradeResponse `json:"nextUpgrade"`
	Effects           []string             `json:"effects"`
	CreatedAt         time.Time            `json:"createdAt"`
	UpdatedAt         time.Time            `json:"updatedAt"`
}

type nextUpgradeResponse struct {
	TargetLevel     int                    `json:"targetLevel"`
	Cost            resourceValuesResponse `json:"cost"`
	DurationSeconds int                    `json:"durationSeconds"`
	CanUpgrade      bool                   `json:"canUpgrade"`
	BlockedReason   *string                `json:"blockedReason"`
}

type buildingsEnvelope struct {
	Buildings []buildingResponse `json:"buildings"`
}

type buildingUpgradeEnvelope struct {
	Building  buildingResponse  `json:"building"`
	Resources resourcesResponse `json:"resources"`
}

func NewBuildingHandler(buildings *service.BuildingService) BuildingHandler {
	return BuildingHandler{buildings: buildings}
}

func (h BuildingHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	buildings, err := h.buildings.Current(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrBuildingKingdomNotFound) {
			return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting buildings")
		}
		return err
	}

	response := make([]buildingResponse, 0, len(buildings))
	for _, building := range buildings {
		response = append(response, newBuildingResponse(building))
	}

	return c.JSON(http.StatusOK, buildingsEnvelope{Buildings: response})
}

func (h BuildingHandler) Upgrade(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	result, err := h.buildings.Upgrade(c.Request().Context(), userID, c.Param("type"))
	if err != nil {
		return h.handleUpgradeError(c, err)
	}

	return c.JSON(http.StatusOK, buildingUpgradeEnvelope{
		Building:  newBuildingResponse(result.Building),
		Resources: newResourcesResponse(result.Resources),
	})
}

func (h BuildingHandler) handleUpgradeError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrBuildingKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting buildings")
	case errors.Is(err, service.ErrInvalidBuildingType):
		return JSONError(c, http.StatusBadRequest, "invalid_building_type", "Building type is invalid")
	case errors.Is(err, service.ErrBuildingNotFound):
		return JSONError(c, http.StatusNotFound, "building_not_found", "Building was not found")
	case errors.Is(err, service.ErrBuildingAlreadyUpgrading):
		return JSONError(c, http.StatusConflict, "building_already_upgrading", "Building is already upgrading")
	case errors.Is(err, service.ErrBuildingMaxLevel):
		return JSONError(c, http.StatusConflict, "building_max_level", "Building is already at max level")
	case errors.Is(err, service.ErrInsufficientResources):
		return JSONError(c, http.StatusConflict, "insufficient_resources", "Not enough resources")
	default:
		return err
	}
}

func newBuildingResponse(view service.BuildingView) buildingResponse {
	building := view.Building
	return buildingResponse{
		ID:                building.ID,
		KingdomID:         building.KingdomID,
		Type:              building.Type,
		Label:             view.Label,
		Level:             building.Level,
		MaxLevel:          view.MaxLevel,
		IsUpgrading:       building.IsUpgrading(),
		UpgradeStartedAt:  building.UpgradeStartedAt,
		UpgradeFinishesAt: building.UpgradeFinishesAt,
		NextUpgrade:       newNextUpgradeResponse(view.Next),
		Effects:           view.Effects,
		CreatedAt:         building.CreatedAt,
		UpdatedAt:         building.UpdatedAt,
	}
}

func newNextUpgradeResponse(next *service.BuildingNextUpgrade) *nextUpgradeResponse {
	if next == nil {
		return nil
	}

	return &nextUpgradeResponse{
		TargetLevel:     next.TargetLevel,
		Cost:            newResourceValuesResponse(next.Cost),
		DurationSeconds: next.DurationSeconds,
		CanUpgrade:      next.CanUpgrade,
		BlockedReason:   next.BlockedReason,
	}
}
