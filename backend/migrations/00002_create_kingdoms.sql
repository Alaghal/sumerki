-- +goose Up
CREATE TABLE kingdoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    culture TEXT NOT NULL,
    patron TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kingdoms_name_length CHECK (char_length(name) BETWEEN 3 AND 32),
    CONSTRAINT kingdoms_culture_valid CHECK (
        culture IN (
            'northern_principality',
            'lizard_grad',
            'free_posad'
        )
    ),
    CONSTRAINT kingdoms_patron_valid CHECK (
        patron IS NULL
        OR patron IN (
            'independent',
            'empire_of_dusk',
            'old_pact'
        )
    )
);

-- +goose Down
DROP TABLE IF EXISTS kingdoms;
