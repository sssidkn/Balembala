-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN app_password_hash;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN app_password_hash VARCHAR(100);
-- +goose StatementEnd