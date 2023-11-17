-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics(
    id VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL,
    delta BIGINT DEFAULT NULL,
    value DOUBLE PRECISION DEFAULT NULL,
    PRIMARY KEY (id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
