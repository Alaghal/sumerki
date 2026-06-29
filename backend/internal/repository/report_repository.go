package repository

import (
	"context"
	"database/sql"

	"sumerki/backend/internal/domain"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateMissionReport(ctx context.Context, kingdomID string, missionID string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte) (domain.MissionReport, error) {
	const query = `
		INSERT INTO mission_reports (kingdom_id, mission_id, title, body, result, rewards_json, losses_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, is_read, created_at
	`

	return scanMissionReport(r.db.QueryRowContext(ctx, query, kingdomID, missionID, title, body, result, rewardsJSON, lossesJSON))
}

func (r *ReportRepository) ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.MissionReport, error) {
	const query = `
		SELECT id::text, kingdom_id::text, mission_id::text, type, title, body, result, rewards_json, losses_json, is_read, created_at
		FROM mission_reports
		WHERE kingdom_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`

	rows, err := r.db.QueryContext(ctx, query, kingdomID)
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
