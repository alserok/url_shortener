-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS urls (
    url text UNIQUE NOT NULL,
    shortened_url varchar(20) NOT NULL PRIMARY KEY
);

CREATE INDEX IF NOT EXISTS shortened_url_hash_index ON urls USING HASH (shortened_url);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
