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

-- name: UpdateProductPrice :exec
UPDATE products SET price = $2 WHERE id = $1;


-- name: GetFullProductDetail :many
SELECT 
  p.id,
  p.product_name,
  i.quantity_on_hand,
  p.price,
  p.reorder_level,
  CASE 
    WHEN i.quantity_on_hand = 0 THEN -2
    WHEN i.quantity_on_hand <= p.reorder_level THEN -1
    WHEN i.quantity_on_hand <= p.reorder_level * 1.5 THEN 0
    ELSE 1
  END AS stock_status
FROM products p
JOIN inventory i
ON p.id = i.product_id;