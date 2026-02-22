-- +goose up
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    price INTEGER NOT NULL,
    reorder_level INTEGER NOT NULL
);

-- +goose down
DROP TABLE products;