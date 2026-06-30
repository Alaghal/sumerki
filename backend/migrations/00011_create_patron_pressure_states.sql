-- +goose Up
CREATE TABLE patron_pressure_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL UNIQUE REFERENCES kingdoms(id) ON DELETE CASCADE,
    patron TEXT NOT NULL,
    tribute_debt_gold BIGINT NOT NULL DEFAULT 0,
    tribute_debt_food BIGINT NOT NULL DEFAULT 0,
    contribution_debt_food BIGINT NOT NULL DEFAULT 0,
    pressure_level INTEGER NOT NULL DEFAULT 0,
    crisis_status TEXT NOT NULL DEFAULT 'none',
    crisis_started_at TIMESTAMPTZ NULL,
    next_tribute_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_resolved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    delay_until TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT patron_pressure_states_patron_valid CHECK (
        patron IN ('independent', 'empire_of_dusk', 'old_pact')
    ),
    CONSTRAINT patron_pressure_states_tribute_debt_gold_nonnegative CHECK (tribute_debt_gold >= 0),
    CONSTRAINT patron_pressure_states_tribute_debt_food_nonnegative CHECK (tribute_debt_food >= 0),
    CONSTRAINT patron_pressure_states_contribution_debt_food_nonnegative CHECK (contribution_debt_food >= 0),
    CONSTRAINT patron_pressure_states_pressure_level_range CHECK (pressure_level BETWEEN 0 AND 100),
    CONSTRAINT patron_pressure_states_crisis_status_valid CHECK (
        crisis_status IN ('none', 'warning', 'active', 'delayed')
    )
);

CREATE INDEX patron_pressure_states_kingdom_id_idx ON patron_pressure_states (kingdom_id);
CREATE INDEX patron_pressure_states_patron_idx ON patron_pressure_states (patron);
CREATE INDEX patron_pressure_states_next_tribute_at_idx ON patron_pressure_states (next_tribute_at);
CREATE INDEX patron_pressure_states_crisis_status_idx ON patron_pressure_states (crisis_status);

INSERT INTO patron_pressure_states (kingdom_id, patron, next_tribute_at)
SELECT
    id,
    patron,
    CASE
        WHEN patron = 'empire_of_dusk' THEN now() + interval '1 hour'
        WHEN patron = 'old_pact' THEN now() + interval '2 hours'
        ELSE now()
    END
FROM kingdoms
WHERE patron IS NOT NULL
ON CONFLICT (kingdom_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS patron_pressure_states;
