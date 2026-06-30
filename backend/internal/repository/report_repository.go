package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"
)

var ErrReportNotFound = errors.New("report not found")

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateMissionReport(ctx context.Context, kingdomID string, missionID string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte, phasesJSON []byte) (domain.MissionReport, error) {
	return r.CreateReport(ctx, kingdomID, &missionID, "pve_mission", title, body, result, rewardsJSON, lossesJSON, phasesJSON)
}

func (r *ReportRepository) CreateReport(ctx context.Context, kingdomID string, missionID *string, reportType string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte, phasesJSON []byte) (domain.MissionReport, error) {
	const query = `
		INSERT INTO mission_reports (kingdom_id, mission_id, type, title, body, result, rewards_json, losses_json, phases_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, phases_json, is_read, created_at
	`

	return scanMissionReport(r.db.QueryRowContext(ctx, query, kingdomID, missionID, reportType, title, body, result, rewardsJSON, lossesJSON, phasesJSON))
}

func (r *ReportRepository) ListByKingdomID(ctx context.Context, kingdomID string, limit int, offset int) ([]domain.MissionReport, error) {
	const query = `
		SELECT id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, phases_json, is_read, created_at
		FROM mission_reports
		WHERE kingdom_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, kingdomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []domain.MissionReport{}
	for rows.Next() {
		report, err := scanMissionReport(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportRepository) CountUnreadByKingdomID(ctx context.Context, kingdomID string) (int64, error) {
	const query = `
		SELECT count(*)
		FROM mission_reports
		WHERE kingdom_id = $1 AND is_read = false
	`

	var count int64
	if err := r.db.QueryRowContext(ctx, query, kingdomID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) FindByIDAndKingdomID(ctx context.Context, reportID string, kingdomID string) (domain.MissionReport, error) {
	const query = `
		SELECT id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, phases_json, is_read, created_at
		FROM mission_reports
		WHERE id = $1 AND kingdom_id = $2
	`

	report, err := scanMissionReport(r.db.QueryRowContext(ctx, query, reportID, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.MissionReport{}, ErrReportNotFound
	}
	if err != nil {
		return domain.MissionReport{}, err
	}
	return report, nil
}

func (r *ReportRepository) MarkRead(ctx context.Context, reportID string, kingdomID string) (domain.MissionReport, error) {
	const query = `
		UPDATE mission_reports
		SET is_read = true
		WHERE id = $1 AND kingdom_id = $2
		RETURNING id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, phases_json, is_read, created_at
	`

	report, err := scanMissionReport(r.db.QueryRowContext(ctx, query, reportID, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.MissionReport{}, ErrReportNotFound
	}
	if err != nil {
		return domain.MissionReport{}, err
	}
	return report, nil
}

func scanMissionReport(row scanner) (domain.MissionReport, error) {
	var report domain.MissionReport
	var missionID sql.NullString
	err := row.Scan(
		&report.ID,
		&report.KingdomID,
		&missionID,
		&report.Type,
		&report.Title,
		&report.Body,
		&report.Result,
		&report.RewardsJSON,
		&report.LossesJSON,
		&report.PhasesJSON,
		&report.IsRead,
		&report.CreatedAt,
	)
	if err != nil {
		return domain.MissionReport{}, err
	}
	if missionID.Valid {
		report.MissionID = &missionID.String
	}

	return report, nil
}
