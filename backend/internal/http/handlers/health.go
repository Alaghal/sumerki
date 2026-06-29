package handlers

import (
	"database/sql"
	"net/http"

	"sumerki/backend/internal/db"

	"github.com/labstack/echo/v4"
)

type Health struct {
	database *sql.DB
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadyResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func NewHealth(database *sql.DB) Health {
	return Health{database: database}
}

func (h Health) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

func (h Health) Ready(c echo.Context) error {
	if err := db.Ping(c.Request().Context(), h.database); err != nil {
		return JSONError(c, http.StatusServiceUnavailable, "database_unavailable", "Database is not reachable")
	}

	return c.JSON(http.StatusOK, ReadyResponse{
		Status:   "ready",
		Database: "ok",
	})
}
