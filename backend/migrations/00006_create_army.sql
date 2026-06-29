-- +goose Up
CREATE TABLE kingdom_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    unit_type TEXT NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kingdom_units_unique_type UNIQUE (kingdom_id, unit_type),
    CONSTRAINT kingdom_units_amount_nonnegative CHECK (amount >= 0),
    CONSTRAINT kingdom_units_type_valid CHECK (unit_type IN (
        'militia',
        'spearmen',
        'archers',
        'cavalry',
        'scouts'
    ))
);

CREATE TABLE unit_training_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    unit_type TEXT NOT NULL,
    amount BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'training',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finishes_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unit_training_orders_amount_positive CHECK (amount > 0),
    CONSTRAINT unit_training_orders_status_valid CHECK (status IN ('training', 'completed')),
    CONSTRAINT unit_training_orders_type_valid CHECK (unit_type IN (
        'militia',
        'spearmen',
        'archers',
        'cavalry',
        'scouts'
    ))
);

CREATE INDEX unit_training_orders_kingdom_id_idx ON unit_training_orders (kingdom_id);
CREATE INDEX unit_training_orders_status_idx ON unit_training_orders (status);
CREATE INDEX unit_training_orders_finishes_at_idx ON unit_training_orders (finishes_at);

INSERT INTO kingdom_units (kingdom_id, unit_type, amount)
SELECT kingdoms.id, units.unit_type, units.amount
FROM kingdoms
CROSS JOIN (
    VALUES
        ('militia', 10),
        ('spearmen', 0),
        ('archers', 0),
        ('cavalry', 0),
        ('scouts', 2)
) AS units(unit_type, amount)
ON CONFLICT (kingdom_id, unit_type) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS unit_training_orders;
DROP TABLE IF EXISTS kingdom_units;
