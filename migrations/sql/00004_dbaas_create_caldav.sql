-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS caldav (
    isu BIGINT PRIMARY KEY,
    ical BYTEA NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS caldav;
-- +goose StatementEnd