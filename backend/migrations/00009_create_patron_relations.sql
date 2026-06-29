-- +goose Up
CREATE TABLE patron_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kingdom_id UUID NOT NULL UNIQUE REFERENCES kingdoms(id) ON DELETE CASCADE,
    patron TEXT NOT NULL,
    favor INTEGER NOT NULL DEFAULT 0,
    standing TEXT NOT NULL DEFAULT 'neutral',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT patron_relations_patron_valid CHECK (
        patron IN (
            'independent',
            'empire_of_dusk',
            'old_pact'
        )
    ),
    CONSTRAINT patron_relations_favor_range CHECK (favor BETWEEN -100 AND 100),
    CONSTRAINT patron_relations_standing_valid CHECK (
        standing IN (
            'hostile',
            'cold',
            'neutral',
            'warm',
            'loyal'
        )
    )
);

INSERT INTO patron_relations (kingdom_id, patron)
SELECT id, patron
FROM kingdoms
WHERE patron IS NOT NULL
ON CONFLICT (kingdom_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS patron_relations;
