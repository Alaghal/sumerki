-- +goose Up
CREATE TABLE missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    mission_key TEXT NOT NULL,
    mission_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finishes_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NULL,
    result_json JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT missions_status_valid CHECK (status IN ('active', 'completed')),
    CONSTRAINT missions_type_valid CHECK (mission_type IN ('expedition', 'scouting'))
);

CREATE INDEX missions_kingdom_id_idx ON missions (kingdom_id);
CREATE INDEX missions_status_idx ON missions (status);
CREATE INDEX missions_finishes_at_idx ON missions (finishes_at);

CREATE TABLE mission_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mission_id UUID NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    unit_type TEXT NOT NULL,
    amount_sent BIGINT NOT NULL,
    amount_lost BIGINT NOT NULL DEFAULT 0,
    amount_returned BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT mission_units_amount_sent_positive CHECK (amount_sent > 0),
    CONSTRAINT mission_units_amount_lost_nonnegative CHECK (amount_lost >= 0),
    CONSTRAINT mission_units_amount_returned_nonnegative CHECK (amount_returned >= 0),
    CONSTRAINT mission_units_amounts_valid CHECK (amount_lost + amount_returned <= amount_sent),
    CONSTRAINT mission_units_type_valid CHECK (unit_type IN (
        'militia',
        'spearmen',
        'archers',
        'cavalry',
        'scouts'
    ))
);

CREATE TABLE mission_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    mission_id UUID NULL REFERENCES missions(id) ON DELETE SET NULL,
    type TEXT NOT NULL DEFAULT 'pve_mission',
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    result TEXT NOT NULL,
    rewards_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    losses_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT mission_reports_type_valid CHECK (type IN ('pve_mission')),
    CONSTRAINT mission_reports_result_valid CHECK (result IN ('success', 'partial_success', 'failure'))
);

CREATE INDEX mission_reports_kingdom_id_idx ON mission_reports (kingdom_id);
CREATE INDEX mission_reports_created_at_idx ON mission_reports (created_at);

-- +goose Down
DROP TABLE IF EXISTS mission_reports;
DROP TABLE IF EXISTS mission_units;
DROP TABLE IF EXISTS missions;
