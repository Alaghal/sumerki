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

type ReportHandler struct {
	missions *service.MissionService
}

type reportsEnvelope struct {
	Reports     []missionReportResponse `json:"reports"`
	Pagination  paginationResponse      `json:"pagination"`
	UnreadCount int64                   `json:"unreadCount"`
}

type reportEnvelope struct {
	Report missionReportResponse `json:"report"`
}

type paginationResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type missionReportResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Phases    []reportPhaseResponse  `json:"phases"`
	Result    string                 `json:"result"`
	Rewards   resourceValuesResponse `json:"rewards"`
	Losses    map[string]int64       `json:"losses"`
	IsRead    bool                   `json:"isRead"`
	CreatedAt time.Time              `json:"createdAt"`
}

type reportPhaseResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewReportHandler(missions *service.MissionService) ReportHandler {
	return ReportHandler{missions: missions}
}

func (h ReportHandler) Me(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	limit := parsePositiveInt(c.QueryParam("limit"), 20)
	if limit > 50 {
		limit = 50
	}
	offset := parsePositiveInt(c.QueryParam("offset"), 0)

	result, err := h.missions.Reports(c.Request().Context(), userID, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrMissionKingdomNotFound) {
			return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting reports")
		}
		return err
	}

	response := make([]missionReportResponse, 0, len(result.Reports))
	for _, report := range result.Reports {
		response = append(response, newMissionReportResponse(report))
	}

	return c.JSON(http.StatusOK, reportsEnvelope{
		Reports:     response,
		Pagination:  paginationResponse{Limit: limit, Offset: offset},
		UnreadCount: result.UnreadCount,
	})
}

func (h ReportHandler) Detail(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	report, err := h.missions.Report(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return reportError(c, err)
	}

	return c.JSON(http.StatusOK, reportEnvelope{Report: newMissionReportResponse(report)})
}

func (h ReportHandler) MarkRead(c echo.Context) error {
	userID, ok := appmiddleware.UserID(c)
	if !ok {
		return JSONError(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	}

	report, err := h.missions.MarkReportRead(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return reportError(c, err)
	}

	return c.JSON(http.StatusOK, reportEnvelope{Report: newMissionReportResponse(report)})
}

func newMissionReportResponse(view service.MissionReportView) missionReportResponse {
	report := view.Report
	phases := make([]reportPhaseResponse, 0, len(view.Phases))
	for _, phase := range view.Phases {
		phases = append(phases, reportPhaseResponse{
			Title: phase.Title,
			Body:  phase.Body,
		})
	}

	return missionReportResponse{
		ID:        report.ID,
		Type:      report.Type,
		Title:     report.Title,
		Body:      report.Body,
		Phases:    phases,
		Result:    report.Result,
		Rewards:   newResourceValuesResponse(view.Rewards),
		Losses:    view.Losses,
		IsRead:    report.IsRead,
		CreatedAt: report.CreatedAt,
	}
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

func reportError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrMissionKingdomNotFound):
		return JSONError(c, http.StatusNotFound, "kingdom_not_found", "Create a kingdom before requesting reports")
	case errors.Is(err, service.ErrReportNotFound):
		return JSONError(c, http.StatusNotFound, "report_not_found", "Report not found")
	default:
		return err
	}
}
