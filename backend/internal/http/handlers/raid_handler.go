package handlers

import (
	"errors"
	"net/http"
	"time"

	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type RaidHandler struct {
	raids *service.RaidService
}

type neighborsEnvelope struct {
	Neighbors []neighborResponse `json:"neighbors"`
}

type neighborResponse struct {
	KingdomID     string  `json:"kingdomId"`
	Name          string  `json:"name"`
	Culture       string  `json:"culture"`
	Patron        *string `json:"patron"`
	Dread         int     `json:"dread"`
	PowerEstimate string  `json:"powerEstimate"`
	CanRaid       bool    `json:"canRaid"`
	BlockedReason *string `json:"blockedReason"`
}

type raidsEnvelope struct {
	Raids []raidResponse `json:"raids"`
}

type startRaidRequest struct {
	DefenderKingdomID string          `json:"defenderKingdomId"`
	Units             []startRaidUnit `json:"units"`
}

type startRaidUnit struct {
	UnitType string `json:"unitType"`
	Amount   int64  `json:"amount"`
}

type startRaidEnvelope struct {
	Raid raidResponse `json:"raid"`
	Army armyResponse `json:"army"`
}

type raidResponse struct {
	ID                  string                 `json:"id"`
	AttackerKingdomID   string                 `json:"attackerKingdomId"`
	AttackerKingdomName string                 `json:"attackerKingdomName"`
	DefenderKingdomID   string                 `json:"defenderKingdomId"`
	DefenderKingdomName string                 `json:"defenderKingdomName"`
	Status              string                 `json:"status"`
	Result              *string                `json:"result"`
	StartedAt           time.Time              `json:"startedAt"`
	ArrivesAt           time.Time              `json:"arrivesAt"`
	CompletedAt         *time.Time             `json:"completedAt"`
	Units               []raidUnitResponse     `json:"units"`
	Loot                resourceValuesResponse `json:"loot"`
}

type raidUnitResponse struct {
	UnitType       string `json:"unitType"`
	UnitLabel      string `json:"unitLabel"`
	AmountSent     int64  `json:"amountSent"`
	AmountLost     int64  `json:"amountLost"`
	AmountReturned int64  `json:"amountReturned"`
}

func NewRaidHandler(raids *service.RaidService) RaidHandler {
	return RaidHandler{raids: raids}
}

func (h RaidHandler) Neighbors(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	neighbors, err := h.raids.Neighbors(c.Request().Context(), userID)
	if err != nil {
		return raidError(c, err)
	}
	response := make([]neighborResponse, 0, len(neighbors))
	for _, neighbor := range neighbors {
		response = append(response, neighborResponse{
			KingdomID:     neighbor.Kingdom.ID,
			Name:          neighbor.Kingdom.Name,
			Culture:       neighbor.Kingdom.Culture,
			Patron:        neighbor.Kingdom.Patron,
			Dread:         neighbor.Kingdom.Dread,
			PowerEstimate: neighbor.PowerEstimate,
			CanRaid:       neighbor.CanRaid,
			BlockedReason: neighbor.BlockedReason,
		})
	}
	return c.JSON(http.StatusOK, neighborsEnvelope{Neighbors: response})
}

func (h RaidHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	raids, err := h.raids.Current(c.Request().Context(), userID)
	if err != nil {
		return raidError(c, err)
	}
	response := make([]raidResponse, 0, len(raids))
	for _, raid := range raids {
		response = append(response, newRaidResponse(raid))
	}
	return c.JSON(http.StatusOK, raidsEnvelope{Raids: response})
}

func (h RaidHandler) Start(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	var req startRaidRequest
	if err := c.Bind(&req); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request", "Invalid JSON request")
	}
	units := make([]service.StartRaidUnit, 0, len(req.Units))
	for _, unit := range req.Units {
		units = append(units, service.StartRaidUnit{UnitType: unit.UnitType, Amount: unit.Amount})
	}

	result, err := h.raids.Start(c.Request().Context(), userID, req.DefenderKingdomID, units)
	if err != nil {
		return raidError(c, err)
	}
	return c.JSON(http.StatusOK, startRaidEnvelope{
		Raid: newRaidResponse(result.Raid),
		Army: newArmyResponse(result.Army),
	})
}

func newRaidResponse(view service.RaidView) raidResponse {
	raid := view.Raid
	units := make([]raidUnitResponse, 0, len(view.Units))
	for _, unit := range view.Units {
		units = append(units, raidUnitResponse{
			UnitType:       unit.Unit.UnitType,
			UnitLabel:      unit.Label,
			AmountSent:     unit.Unit.AmountSent,
			AmountLost:     unit.Unit.AmountLost,
			AmountReturned: unit.Unit.AmountReturned,
		})
	}
	return raidResponse{
		ID:                  raid.ID,
		AttackerKingdomID:   raid.AttackerKingdomID,
		AttackerKingdomName: view.AttackerKingdomName,
		DefenderKingdomID:   raid.DefenderKingdomID,
		DefenderKingdomName: view.DefenderKingdomName,
		Status:              raid.Status,
		Result:              raid.Result,
		StartedAt:           raid.StartedAt,
		ArrivesAt:           raid.ArrivesAt,
		CompletedAt:         raid.CompletedAt,
		Units:               units,
		Loot:                newResourceValuesResponse(view.Loot),
	}
}

func raidError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrRaidKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting raids")
	case errors.Is(err, service.ErrRaidTargetNotFound):
		return JSONError(c, http.StatusNotFound, "target_not_found", "Raid target not found")
	case errors.Is(err, service.ErrCannotRaidSelf):
		return JSONError(c, http.StatusBadRequest, "cannot_raid_self", "Cannot raid your own kingdom")
	case errors.Is(err, service.ErrInvalidRaidUnitType):
		return JSONError(c, http.StatusBadRequest, "invalid_unit_type", "Unit type is invalid")
	case errors.Is(err, service.ErrInvalidRaidUnitAmount):
		return JSONError(c, http.StatusBadRequest, "invalid_unit_amount", "Unit amount is invalid")
	case errors.Is(err, service.ErrRaidInsufficientUnits):
		return JSONError(c, http.StatusBadRequest, "insufficient_units", "Not enough units")
	case errors.Is(err, service.ErrRaidRequirementsNotMet):
		return JSONError(c, http.StatusBadRequest, "raid_requirements_not_met", "Raid requirements are not met")
	case errors.Is(err, service.ErrTargetNewbieProtected):
		return JSONError(c, http.StatusBadRequest, "target_newbie_protected", "Target is protected from raids")
	case errors.Is(err, service.ErrTargetTooWeak):
		return JSONError(c, http.StatusBadRequest, "target_too_weak", "Target is too weak to raid")
	case errors.Is(err, service.ErrRaidCooldownActive):
		return JSONError(c, http.StatusBadRequest, "raid_cooldown_active", "This target was raided too recently")
	case errors.Is(err, service.ErrTargetUnderProtection):
		return JSONError(c, http.StatusBadRequest, "target_under_protection", "Target is under temporary raid protection")
	default:
		return err
	}
}
