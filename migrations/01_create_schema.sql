-- +goose Up
CREATE SCHEMA gophkeeper AUTHORIZATION postgres;
GRANT USAGE ON SCHEMA gophkeeper TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gophkeeper TO postgres;

-- +goose Down
DROP SCHEMA gophkeeper;