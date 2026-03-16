-- +goose up
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    product_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    price DECIMAL (10, 2) NOT NULL,
    reorder_level INTEGER NOT NULL DEFAULT 5
);

CREATE TABLE IF NOT EXISTS purchases (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    quantity_added INTEGER NOT NULL CHECK (quantity_added > 0),
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);


CREATE TABLE IF NOT EXISTS inventory (
    product_id TEXT PRIMARY KEY,
    quantity_on_hand INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);


CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    total_amount DECIMAL(10,2) NOT NULL,
    sale_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales_items (
    id SERIAL PRIMARY KEY,
    sales_id INTEGER NOT NULL,
    product_id TEXT NOT NULL,
    quantity_sold INTEGER NOT NULL CHECK (quantity_sold > 0),
    price_at_sale DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (sales_id) REFERENCES sales(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS adjustments (
    id SERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    quantity_changed INTEGER NOT NULL,
    reason TEXT NOT NULL,
    adjustment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +goose down
DROP TABLE adjustments;
DROP TABLE sales_items;
DROP TABLE purchases;
DROP TABLE sales;
DROP TABLE inventory;
DROP TABLE products;