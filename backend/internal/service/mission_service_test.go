package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

func TestReportsReturnsUnreadCountAndEmptyLegacyPhases(t *testing.T) {
	ctx := context.Background()
	reports := &fakeMissionReportRepository{
		reports: []domain.MissionReport{
			{
				ID:          "report-1",
				KingdomID:   "kingdom-1",
				Type:        "pve_mission",
				Title:       "Old report",
				Body:        "Legacy body",
				Result:      "success",
				RewardsJSON: []byte(`{"gold":1}`),
				LossesJSON:  []byte(`{}`),
				IsRead:      false,
				CreatedAt:   time.Now(),
			},
		},
	}
	service := NewMissionService(
		&fakeMissionKingdomRepository{kingdom: domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}},
		&fakeMissionRepository{},
		reports,
		nil,
		nil,
	)

	result, err := service.Reports(ctx, "user-1", 20, 0)
	if err != nil {
		t.Fatalf("Reports returned error: %v", err)
	}
	if result.UnreadCount != 1 {
		t.Fatalf("UnreadCount = %d, want 1", result.UnreadCount)
	}
	if len(result.Reports) != 1 {
		t.Fatalf("len(Reports) = %d, want 1", len(result.Reports))
	}
	if len(result.Reports[0].Phases) != 0 {
		t.Fatalf("len(Phases) = %d, want 0 for legacy report", len(result.Reports[0].Phases))
	}
}

func TestReportRejectsAnotherKingdomReport(t *testing.T) {
	ctx := context.Background()
	service := NewMissionService(
		&fakeMissionKingdomRepository{kingdom: domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}},
		&fakeMissionRepository{},
		&fakeMissionReportRepository{
			reports: []domain.MissionReport{
				{ID: "report-1", KingdomID: "kingdom-2", Type: "pve_mission", Result: "success"},
			},
		},
		nil,
		nil,
	)

	if _, err := service.Report(ctx, "user-1", "report-1"); !errors.Is(err, ErrReportNotFound) {
		t.Fatalf("Report error = %v, want ErrReportNotFound", err)
	}
}

func TestMarkReportReadIsIdempotent(t *testing.T) {
	ctx := context.Background()
	reports := &fakeMissionReportRepository{
		reports: []domain.MissionReport{
			{ID: "report-1", KingdomID: "kingdom-1", Type: "pve_mission", Result: "success", IsRead: false},
		},
	}
	service := NewMissionService(
		&fakeMissionKingdomRepository{kingdom: domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}},
		&fakeMissionRepository{},
		reports,
		nil,
		nil,
	)

	first, err := service.MarkReportRead(ctx, "user-1", "report-1")
	if err != nil {
		t.Fatalf("first MarkReportRead returned error: %v", err)
	}
	second, err := service.MarkReportRead(ctx, "user-1", "report-1")
	if err != nil {
		t.Fatalf("second MarkReportRead returned error: %v", err)
	}
	if !first.Report.IsRead || !second.Report.IsRead {
		t.Fatalf("MarkReportRead did not keep report read")
	}
}

func TestMissionReportTemplateIncludesPhases(t *testing.T) {
	template := gameconfig.MissionReportTemplate("black_forest_expedition", "success")
	if template.Title == "" || template.Body == "" {
		t.Fatalf("template title/body must not be empty")
	}
	if len(template.Phases) == 0 {
		t.Fatalf("len(Phases) = 0, want narrative phases")
	}
}

type fakeMissionKingdomRepository struct {
	kingdom domain.Kingdom
}

func (r *fakeMissionKingdomRepository) Create(ctx context.Context, userID string, name string, culture string) (domain.Kingdom, error) {
	return domain.Kingdom{}, nil
}

func (r *fakeMissionKingdomRepository) FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error) {
	if r.kingdom.UserID != userID {
		return domain.Kingdom{}, repository.ErrKingdomNotFound
	}
	return r.kingdom, nil
}

type fakeMissionRepository struct{}

func (r *fakeMissionRepository) CreateMission(ctx context.Context, kingdomID string, missionKey string, missionType string, startedAt time.Time, finishesAt time.Time) (domain.Mission, error) {
	return domain.Mission{}, nil
}

func (r *fakeMissionRepository) CreateMissionUnit(ctx context.Context, missionID string, unitType string, amount int64) (domain.MissionUnit, error) {
	return domain.MissionUnit{}, nil
}

func (r *fakeMissionRepository) ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error) {
	return nil, nil
}

func (r *fakeMissionRepository) ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error) {
	return nil, nil
}

func (r *fakeMissionRepository) ListUnitsByMissionID(ctx context.Context, missionID string) ([]domain.MissionUnit, error) {
	return nil, nil
}

func (r *fakeMissionRepository) UpdateMissionUnitResult(ctx context.Context, missionUnitID string, lost int64, returned int64) error {
	return nil
}

func (r *fakeMissionRepository) CompleteMission(ctx context.Context, missionID string, completedAt time.Time, resultJSON []byte) error {
	return nil
}

type fakeMissionReportRepository struct {
	reports []domain.MissionReport
}

func (r *fakeMissionReportRepository) CreateMissionReport(ctx context.Context, kingdomID string, missionID string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte, phasesJSON []byte) (domain.MissionReport, error) {
	report := domain.MissionReport{
		ID:          missionID + "-report",
		KingdomID:   kingdomID,
		MissionID:   &missionID,
		Type:        "pve_mission",
		Title:       title,
		Body:        body,
		Result:      result,
		RewardsJSON: rewardsJSON,
		LossesJSON:  lossesJSON,
		PhasesJSON:  phasesJSON,
	}
	r.reports = append(r.reports, report)
	return report, nil
}

func (r *fakeMissionReportRepository) ListByKingdomID(ctx context.Context, kingdomID string, limit int, offset int) ([]domain.MissionReport, error) {
	filtered := []domain.MissionReport{}
	for _, report := range r.reports {
		if report.KingdomID == kingdomID {
			filtered = append(filtered, report)
		}
	}
	if offset >= len(filtered) {
		return []domain.MissionReport{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

func (r *fakeMissionReportRepository) CountUnreadByKingdomID(ctx context.Context, kingdomID string) (int64, error) {
	var count int64
	for _, report := range r.reports {
		if report.KingdomID == kingdomID && !report.IsRead {
			count++
		}
	}
	return count, nil
}

func (r *fakeMissionReportRepository) FindByIDAndKingdomID(ctx context.Context, reportID string, kingdomID string) (domain.MissionReport, error) {
	for _, report := range r.reports {
		if report.ID == reportID && report.KingdomID == kingdomID {
			return report, nil
		}
	}
	return domain.MissionReport{}, repository.ErrReportNotFound
}

func (r *fakeMissionReportRepository) MarkRead(ctx context.Context, reportID string, kingdomID string) (domain.MissionReport, error) {
	for index, report := range r.reports {
		if report.ID == reportID && report.KingdomID == kingdomID {
			r.reports[index].IsRead = true
			return r.reports[index], nil
		}
	}
	return domain.MissionReport{}, repository.ErrReportNotFound
}
