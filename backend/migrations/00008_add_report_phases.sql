-- +goose Up
ALTER TABLE mission_reports
ADD COLUMN phases_json JSONB NOT NULL DEFAULT '[]'::jsonb;

-- +goose Down
ALTER TABLE mission_reports
DROP COLUMN phases_json;
