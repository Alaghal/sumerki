-- +goose Up
CREATE TABLE rulers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL UNIQUE REFERENCES kingdoms(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    culture TEXT NOT NULL,
    authority INTEGER NOT NULL,
    courage INTEGER NOT NULL,
    cunning INTEGER NOT NULL,
    honor INTEGER NOT NULL,
    cruelty INTEGER NOT NULL,
    ambition INTEGER NOT NULL,
    paranoia INTEGER NOT NULL,
    health_status TEXT NOT NULL DEFAULT 'healthy',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT rulers_name_not_empty CHECK (char_length(name) > 0),
    CONSTRAINT rulers_age_range CHECK (age BETWEEN 18 AND 80),
    CONSTRAINT rulers_culture_valid CHECK (
        culture IN (
            'northern_principality',
            'lizard_grad',
            'free_posad'
        )
    ),
    CONSTRAINT rulers_authority_range CHECK (authority BETWEEN 1 AND 100),
    CONSTRAINT rulers_courage_range CHECK (courage BETWEEN 1 AND 100),
    CONSTRAINT rulers_cunning_range CHECK (cunning BETWEEN 1 AND 100),
    CONSTRAINT rulers_honor_range CHECK (honor BETWEEN 1 AND 100),
    CONSTRAINT rulers_cruelty_range CHECK (cruelty BETWEEN 1 AND 100),
    CONSTRAINT rulers_ambition_range CHECK (ambition BETWEEN 1 AND 100),
    CONSTRAINT rulers_paranoia_range CHECK (paranoia BETWEEN 1 AND 100),
    CONSTRAINT rulers_health_status_valid CHECK (
        health_status IN (
            'healthy',
            'wounded',
            'sick'
        )
    )
);

INSERT INTO rulers (
    kingdom_id,
    name,
    age,
    culture,
    authority,
    courage,
    cunning,
    honor,
    cruelty,
    ambition,
    paranoia,
    health_status
)
SELECT
    kingdoms.id,
    CASE kingdoms.culture
        WHEN 'northern_principality' THEN 'Боривой'
        WHEN 'lizard_grad' THEN 'Шессар'
        ELSE 'Берест'
    END,
    42,
    kingdoms.culture,
    55,
    55,
    55,
    55,
    35,
    55,
    35,
    'healthy'
FROM kingdoms
WHERE NOT EXISTS (
    SELECT 1
    FROM rulers
    WHERE rulers.kingdom_id = kingdoms.id
);

-- +goose Down
DROP TABLE IF EXISTS rulers;
