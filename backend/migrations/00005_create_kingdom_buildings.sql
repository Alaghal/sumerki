-- +goose Up
CREATE TABLE kingdom_buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL REFERENCES kingdoms(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 0,
    upgrade_started_at TIMESTAMPTZ NULL,
    upgrade_finishes_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kingdom_buildings_unique_type UNIQUE (kingdom_id, type),
    CONSTRAINT kingdom_buildings_level_range CHECK (level BETWEEN 0 AND 5),
    CONSTRAINT kingdom_buildings_type_valid CHECK (
        type IN (
            'town_hall',
            'farm',
            'lumberyard',
            'quarry',
            'barracks',
            'market',
            'walls',
            'shrine'
        )
    ),
    CONSTRAINT kingdom_buildings_upgrade_times_valid CHECK (
        (upgrade_started_at IS NULL AND upgrade_finishes_at IS NULL)
        OR (upgrade_started_at IS NOT NULL AND upgrade_finishes_at IS NOT NULL)
    )
);

INSERT INTO kingdom_buildings (kingdom_id, type, level)
SELECT kingdoms.id, buildings.type, buildings.level
FROM kingdoms
CROSS JOIN (
    VALUES
        ('town_hall', 1),
        ('farm', 1),
        ('lumberyard', 1),
        ('quarry', 1),
        ('market', 1),
        ('barracks', 0),
        ('walls', 0),
        ('shrine', 0)
) AS buildings(type, level)
WHERE NOT EXISTS (
    SELECT 1
    FROM kingdom_buildings
    WHERE kingdom_buildings.kingdom_id = kingdoms.id
      AND kingdom_buildings.type = buildings.type
);

-- +goose Down
DROP TABLE IF EXISTS kingdom_buildings;
