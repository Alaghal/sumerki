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

type RulerHandler struct {
	rulers *service.RulerService
}

type rulerResponse struct {
	ID           string    `json:"id"`
	KingdomID    string    `json:"kingdomId"`
	Name         string    `json:"name"`
	Age          int       `json:"age"`
	Culture      string    `json:"culture"`
	Authority    int       `json:"authority"`
	Courage      int       `json:"courage"`
	Cunning      int       `json:"cunning"`
	Honor        int       `json:"honor"`
	Cruelty      int       `json:"cruelty"`
	Ambition     int       `json:"ambition"`
	Paranoia     int       `json:"paranoia"`
	HealthStatus string    `json:"healthStatus"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type rulerEnvelope struct {
	Ruler rulerResponse `json:"ruler"`
}

func NewRulerHandler(rulers *service.RulerService) RulerHandler {
	return RulerHandler{rulers: rulers}
}

func (h RulerHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	ruler, err := h.rulers.Current(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrRulerKingdomNotFound) {
			return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting a ruler")
		}
		return err
	}

	return c.JSON(http.StatusOK, rulerEnvelope{Ruler: newRulerResponse(ruler)})
}

func newRulerResponse(ruler domain.Ruler) rulerResponse {
	return rulerResponse{
		ID:           ruler.ID,
		KingdomID:    ruler.KingdomID,
		Name:         ruler.Name,
		Age:          ruler.Age,
		Culture:      ruler.Culture,
		Authority:    ruler.Authority,
		Courage:      ruler.Courage,
		Cunning:      ruler.Cunning,
		Honor:        ruler.Honor,
		Cruelty:      ruler.Cruelty,
		Ambition:     ruler.Ambition,
		Paranoia:     ruler.Paranoia,
		HealthStatus: ruler.HealthStatus,
		CreatedAt:    ruler.CreatedAt,
		UpdatedAt:    ruler.UpdatedAt,
	}
}
