-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS templates
(
    id      SERIAL PRIMARY KEY,
    title   VARCHAR(255),
    user_id INTEGER NOT NULL,
    message TEXT
);

CREATE TABLE IF NOT EXISTS contacts
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER             NOT NULL,
    name    VARCHAR(255),
    email   VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts_templates
(
    id          SERIAL PRIMARY KEY,
    contact_id  INTEGER NOT NULL,
    template_id INTEGER NOT NULL,
    FOREIGN KEY (contact_id) REFERENCES contacts (id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES templates (id) ON DELETE CASCADE,
    UNIQUE (contact_id, template_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contacts_templates;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd