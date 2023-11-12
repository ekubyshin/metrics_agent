-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics (
    name VARCHAR(20) NOT NULL,
    type VARCHAR(10) NOT NULL,
    delta INT DEFAULT NULL,
    value DOUBLE DEFAULT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
