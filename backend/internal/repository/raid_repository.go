package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
)

var ErrRaidNotFound = errors.New("raid not found")

type RaidRepository struct {
	db *sql.DB
}

func NewRaidRepository(db *sql.DB) *RaidRepository {
	return &RaidRepository{db: db}
}

func (r *RaidRepository) CreateRaid(ctx context.Context, attackerKingdomID string, defenderKingdomID string, startedAt time.Time, arrivesAt time.Time) (domain.Raid, error) {
	const query = `
		INSERT INTO raids (attacker_kingdom_id, defender_kingdom_id, started_at, arrives_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, attacker_kingdom_id::text, defender_kingdom_id::text, status, started_at, arrives_at, completed_at, result, loot_json, attacker_losses_json, defender_losses_json, result_json, created_at, updated_at
	`

	return scanRaid(r.db.QueryRowContext(ctx, query, attackerKingdomID, defenderKingdomID, startedAt, arrivesAt))
}

func (r *RaidRepository) CreateRaidUnit(ctx context.Context, raidID string, side string, unitType string, amount int64) (domain.RaidUnit, error) {
	const query = `
		INSERT INTO raid_units (raid_id, side, unit_type, amount_sent)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, raid_id::text, side, unit_type, amount_sent, amount_lost, amount_returned, created_at, updated_at
	`

	return scanRaidUnit(r.db.QueryRowContext(ctx, query, raidID, side, unitType, amount))
}

func (r *RaidRepository) ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Raid, error) {
	const query = `
		SELECT id::text, attacker_kingdom_id::text, defender_kingdom_id::text, status, started_at, arrives_at, completed_at, result, loot_json, attacker_losses_json, defender_losses_json, result_json, created_at, updated_at
		FROM raids
		WHERE attacker_kingdom_id = $1 OR defender_kingdom_id = $1
		ORDER BY created_at DESC
	`

	return r.list(ctx, query, kingdomID)
}

func (r *RaidRepository) ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.Raid, error) {
	const query = `
		SELECT id::text, attacker_kingdom_id::text, defender_kingdom_id::text, status, started_at, arrives_at, completed_at, result, loot_json, attacker_losses_json, defender_losses_json, result_json, created_at, updated_at
		FROM raids
		WHERE status = 'active'
		  AND (attacker_kingdom_id = $1 OR defender_kingdom_id = $1)
		ORDER BY arrives_at ASC
	`

	return r.list(ctx, query, kingdomID)
}

func (r *RaidRepository) ListUnitsByRaidID(ctx context.Context, raidID string) ([]domain.RaidUnit, error) {
	const query = `
		SELECT id::text, raid_id::text, side, unit_type, amount_sent, amount_lost, amount_returned, created_at, updated_at
		FROM raid_units
		WHERE raid_id = $1
		ORDER BY side, array_position(ARRAY[
			'militia',
			'spearmen',
			'archers',
			'cavalry',
			'scouts'
		]::text[], unit_type)
	`

	rows, err := r.db.QueryContext(ctx, query, raidID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := []domain.RaidUnit{}
	for rows.Next() {
		unit, err := scanRaidUnit(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return units, nil
}

func (r *RaidRepository) UpdateRaidUnitResult(ctx context.Context, raidUnitID string, lost int64, returned int64) error {
	const query = `
		UPDATE raid_units
		SET amount_lost = $2,
			amount_returned = $3,
			updated_at = now()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, raidUnitID, lost, returned)
	return err
}

func (r *RaidRepository) CompleteRaid(ctx context.Context, raidID string, completedAt time.Time, result string, lootJSON []byte, attackerLossesJSON []byte, defenderLossesJSON []byte, resultJSON []byte) error {
	const query = `
		UPDATE raids
		SET status = 'completed',
			completed_at = $2,
			result = $3,
			loot_json = $4,
			attacker_losses_json = $5,
			defender_losses_json = $6,
			result_json = $7,
			updated_at = now()
		WHERE id = $1 AND status = 'active'
	`
	resultExec, err := r.db.ExecContext(ctx, query, raidID, completedAt, result, lootJSON, attackerLossesJSON, defenderLossesJSON, resultJSON)
	if err != nil {
		return err
	}
	rowsAffected, err := resultExec.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRaidNotFound
	}
	return nil
}

func (r *RaidRepository) CountRecentBetween(ctx context.Context, attackerKingdomID string, defenderKingdomID string, since time.Time) (int64, error) {
	const query = `
		SELECT count(*)
		FROM raids
		WHERE attacker_kingdom_id = $1
		  AND defender_kingdom_id = $2
		  AND started_at >= $3
	`
	var count int64
	if err := r.db.QueryRowContext(ctx, query, attackerKingdomID, defenderKingdomID, since).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RaidRepository) CountRecentAgainst(ctx context.Context, defenderKingdomID string, since time.Time) (int64, error) {
	const query = `
		SELECT count(*)
		FROM raids
		WHERE defender_kingdom_id = $1
		  AND started_at >= $2
	`
	var count int64
	if err := r.db.QueryRowContext(ctx, query, defenderKingdomID, since).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RaidRepository) list(ctx context.Context, query string, kingdomID string) ([]domain.Raid, error) {
	rows, err := r.db.QueryContext(ctx, query, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	raids := []domain.Raid{}
	for rows.Next() {
		raid, err := scanRaid(rows)
		if err != nil {
			return nil, err
		}
		raids = append(raids, raid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return raids, nil
}

func scanRaid(row scanner) (domain.Raid, error) {
	var raid domain.Raid
	var completedAt sql.NullTime
	var result sql.NullString
	var lootJSON []byte
	var attackerLossesJSON []byte
	var defenderLossesJSON []byte
	var resultJSON []byte
	err := row.Scan(
		&raid.ID,
		&raid.AttackerKingdomID,
		&raid.DefenderKingdomID,
		&raid.Status,
		&raid.StartedAt,
		&raid.ArrivesAt,
		&completedAt,
		&result,
		&lootJSON,
		&attackerLossesJSON,
		&defenderLossesJSON,
		&resultJSON,
		&raid.CreatedAt,
		&raid.UpdatedAt,
	)
	if err != nil {
		return domain.Raid{}, err
	}
	if completedAt.Valid {
		raid.CompletedAt = &completedAt.Time
	}
	if result.Valid {
		raid.Result = &result.String
	}
	raid.LootJSON = lootJSON
	raid.AttackerLossesJSON = attackerLossesJSON
	raid.DefenderLossesJSON = defenderLossesJSON
	raid.ResultJSON = resultJSON
	return raid, nil
}

func scanRaidUnit(row scanner) (domain.RaidUnit, error) {
	var unit domain.RaidUnit
	err := row.Scan(
		&unit.ID,
		&unit.RaidID,
		&unit.Side,
		&unit.UnitType,
		&unit.AmountSent,
		&unit.AmountLost,
		&unit.AmountReturned,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		return domain.RaidUnit{}, err
	}
	return unit, nil
}
