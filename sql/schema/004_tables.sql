-- +goose up
CREATE TABLE IF NOT EXISTS ewma (
    id BIGSERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ewma DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_ewma
ON ewma(product_id, recorded_at DESC);

ALTER TABLE sales_items
DROP CONSTRAINT sales_items_sales_id_fkey;

ALTER TABLE sales_items
ADD CONSTRAINT sales_items_sales_id_fkey
FOREIGN KEY (sales_id)
REFERENCES sales(id)
ON DELETE CASCADE;

-- +goose down
DROP TABLE 
ewma;