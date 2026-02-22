-- name: CreateProduct :one
INSERT INTO products (id, created_at, updated_at, price, reorder_level, product_name)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;