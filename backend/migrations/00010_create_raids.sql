-- +goose Up
ALTER TABLE kingdoms
ADD COLUMN dread INTEGER NOT NULL DEFAULT 0,
ADD COLUMN honor INTEGER NOT NULL DEFAULT 0;

ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_type_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_type_valid CHECK (
    type IN (
        'pve_mission',
        'pvp_raid_attacker',
        'pvp_raid_defender'
    )
);

ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_result_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_result_valid CHECK (
    result IN (
        'success',
        'partial_success',
        'failure',
        'attacker_success',
        'defender_success',
        'bloody_stalemate',
        'repelled_by_protection'
    )
);

CREATE TABLE raids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    defender_kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    arrives_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NULL,
    result TEXT NULL,
    loot_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    attacker_losses_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    defender_losses_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    result_json JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT raids_different_kingdoms CHECK (attacker_kingdom_id <> defender_kingdom_id),
    CONSTRAINT raids_status_valid CHECK (status IN ('active', 'completed')),
    CONSTRAINT raids_result_valid CHECK (
        result IS NULL
        OR result IN (
            'attacker_success',
            'defender_success',
            'bloody_stalemate',
            'repelled_by_protection'
        )
    )
);

CREATE INDEX raids_attacker_kingdom_id_idx ON raids (attacker_kingdom_id);
CREATE INDEX raids_defender_kingdom_id_idx ON raids (defender_kingdom_id);
CREATE INDEX raids_status_idx ON raids (status);
CREATE INDEX raids_arrives_at_idx ON raids (arrives_at);
CREATE INDEX raids_created_at_idx ON raids (created_at);

CREATE TABLE raid_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    raid_id UUID NOT NULL REFERENCES raids(id) ON DELETE CASCADE,
    side TEXT NOT NULL,
    unit_type TEXT NOT NULL,
    amount_sent BIGINT NOT NULL,
    amount_lost BIGINT NOT NULL DEFAULT 0,
    amount_returned BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT raid_units_side_valid CHECK (side IN ('attacker', 'defender')),
    CONSTRAINT raid_units_amount_sent_nonnegative CHECK (amount_sent >= 0),
    CONSTRAINT raid_units_amount_lost_nonnegative CHECK (amount_lost >= 0),
    CONSTRAINT raid_units_amount_returned_nonnegative CHECK (amount_returned >= 0),
    CONSTRAINT raid_units_amounts_valid CHECK (amount_lost + amount_returned <= amount_sent),
    CONSTRAINT raid_units_type_valid CHECK (unit_type IN (
        'militia',
        'spearmen',
        'archers',
        'cavalry',
        'scouts'
    ))
);

CREATE INDEX raid_units_raid_id_idx ON raid_units (raid_id);

-- +goose Down
DROP TABLE IF EXISTS raid_units;
DROP TABLE IF EXISTS raids;

ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_result_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_result_valid CHECK (result IN ('success', 'partial_success', 'failure'));

ALTER TABLE mission_reports
DROP CONSTRAINT mission_reports_type_valid;

ALTER TABLE mission_reports
ADD CONSTRAINT mission_reports_type_valid CHECK (type IN ('pve_mission'));

ALTER TABLE kingdoms
DROP COLUMN honor,
DROP COLUMN dread;
