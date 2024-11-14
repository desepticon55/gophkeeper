-- +goose Up
CREATE TABLE gophkeeper.secret
(
    name     VARCHAR(255),
    username VARCHAR(255),
    content  BYTEA       NOT NULL,
    type     VARCHAR(30) NOT NULL,
    opt_lock BIGINT      NOT NULL,
    PRIMARY KEY (username, name)
);

-- +goose Down
DROP TABLE gophkeeper.secret;
