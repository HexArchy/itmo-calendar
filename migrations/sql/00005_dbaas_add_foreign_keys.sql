-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_tokens ADD CONSTRAINT fk_user_tokens_users FOREIGN KEY (isu) REFERENCES users(isu) ON DELETE CASCADE;
ALTER TABLE caldav ADD CONSTRAINT fk_caldav_users FOREIGN KEY (isu) REFERENCES users(isu) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE caldav DROP CONSTRAINT IF EXISTS fk_caldav_users;
ALTER TABLE user_tokens DROP CONSTRAINT IF EXISTS fk_user_tokens_users;
-- +goose StatementEnd
