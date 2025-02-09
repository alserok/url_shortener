-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls
(
    url           text UNIQUE        NOT NULL,
    shortened_url varchar(20) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS shortened_url_hash_index ON urls USING HASH (shortened_url);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
DROP INDEX IF EXISTS shortened_url_hash_index;
-- +goose StatementEnd
