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

type ResourcesHandler struct {
	resources *service.ResourcesService
}

type resourceValuesResponse struct {
	Gold       int64 `json:"gold"`
	Food       int64 `json:"food"`
	Wood       int64 `json:"wood"`
	Stone      int64 `json:"stone"`
	Population int64 `json:"population"`
}

type resourcesResponse struct {
	KingdomID         string                 `json:"kingdomId"`
	Gold              int64                  `json:"gold"`
	Food              int64                  `json:"food"`
	Wood              int64                  `json:"wood"`
	Stone             int64                  `json:"stone"`
	Population        int64                  `json:"population"`
	ProductionPerHour resourceValuesResponse `json:"productionPerHour"`
	LastCalculatedAt  time.Time              `json:"lastCalculatedAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

type resourcesEnvelope struct {
	Resources resourcesResponse `json:"resources"`
}

func NewResourcesHandler(resources *service.ResourcesService) ResourcesHandler {
	return ResourcesHandler{resources: resources}
}

func (h ResourcesHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	result, err := h.resources.Current(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrResourcesKingdomNotFound) {
			return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting resources")
		}
		return err
	}

	return c.JSON(http.StatusOK, resourcesEnvelope{Resources: newResourcesResponse(result)})
}

func newResourcesResponse(result service.ResourcesResult) resourcesResponse {
	resources := result.Resources
	return resourcesResponse{
		KingdomID:         resources.KingdomID,
		Gold:              resources.Gold,
		Food:              resources.Food,
		Wood:              resources.Wood,
		Stone:             resources.Stone,
		Population:        resources.Population,
		ProductionPerHour: newResourceValuesResponse(result.ProductionPerHour),
		LastCalculatedAt:  resources.LastCalculatedAt,
		UpdatedAt:         resources.UpdatedAt,
	}
}

func newResourceValuesResponse(values gameconfig.ResourceValues) resourceValuesResponse {
	return resourceValuesResponse{
		Gold:       values.Gold,
		Food:       values.Food,
		Wood:       values.Wood,
		Stone:      values.Stone,
		Population: values.Population,
	}
}
