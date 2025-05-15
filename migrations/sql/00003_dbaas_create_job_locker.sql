-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS job_locks (
    job_name TEXT PRIMARY KEY,
    locked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS job_locks;
-- +goose StatementEnd