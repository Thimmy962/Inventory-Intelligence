-- name: CreateProduct :one
INSERT INTO products (id, created_at, updated_at, product_name, reorder_level, price)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;


-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;