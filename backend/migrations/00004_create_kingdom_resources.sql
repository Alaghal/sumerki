-- +goose Up
CREATE TABLE kingdom_resources (
    kingdom_id UUID PRIMARY KEY REFERENCES kingdoms(id) ON DELETE CASCADE,
    gold BIGINT NOT NULL DEFAULT 500,
    food BIGINT NOT NULL DEFAULT 300,
    wood BIGINT NOT NULL DEFAULT 300,
    stone BIGINT NOT NULL DEFAULT 200,
    population BIGINT NOT NULL DEFAULT 100,
    last_calculated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kingdom_resources_gold_nonnegative CHECK (gold >= 0),
    CONSTRAINT kingdom_resources_food_nonnegative CHECK (food >= 0),
    CONSTRAINT kingdom_resources_wood_nonnegative CHECK (wood >= 0),
    CONSTRAINT kingdom_resources_stone_nonnegative CHECK (stone >= 0),
    CONSTRAINT kingdom_resources_population_nonnegative CHECK (population >= 0)
);

INSERT INTO kingdom_resources (kingdom_id)
SELECT kingdoms.id
FROM kingdoms
WHERE NOT EXISTS (
    SELECT 1
    FROM kingdom_resources
    WHERE kingdom_resources.kingdom_id = kingdoms.id
);

-- +goose Down
DROP TABLE IF EXISTS kingdom_resources;
