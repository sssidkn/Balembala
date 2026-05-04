-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username VARCHAR (100) NOT NULL,
    email VARCHAR (150) NOT NULL UNIQUE,
    password_hash VARCHAR (100) NOT NULL,
    app_password_hash VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd