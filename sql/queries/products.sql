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



-- name: GetFullProductDetailView :many
SELECT * FROM productdetail;


-- name: GetOneFullProductDetail :one
SELECT 
  p.id,
  p.product_name,
  i.quantity_on_hand,
  p.price,
  p.reorder_level,
  p.updated_at,
  i.last_updated,
  CASE 
    WHEN i.quantity_on_hand = 0 THEN -2
    WHEN i.quantity_on_hand <= p.reorder_level THEN -1
    WHEN i.quantity_on_hand <= p.reorder_level * 1.5 THEN 0
    ELSE 1
  END AS stock_status
FROM products p
JOIN inventory i
ON p.id = i.product_id
WHERE p.id = $1;


-- name: EditOneProduct :exec
UPDATE products
SET product_name = $2, price = $3, reorder_level = $4, updated_at = NOW()
WHERE id = $1;


-- name: GetTopProduct :many
SELECT 
  p.id,
  p.product_name,
  i.quantity_on_hand,
  p.price,
  p.reorder_level,
  p.updated_at,
  i.last_updated,
  e.ewma,
  CASE 
    WHEN i.quantity_on_hand = 0 THEN -2
    WHEN i.quantity_on_hand <= p.reorder_level THEN -1
    WHEN i.quantity_on_hand <= p.reorder_level * 1.5 THEN 0
    ELSE 1
  END AS stock_status
FROM products p
JOIN inventory i
  ON p.id = i.product_id
JOIN (
    SELECT *
    FROM (
        SELECT *,
               ROW_NUMBER() OVER (
                   PARTITION BY product_id
                   ORDER BY recorded_at DESC
               ) AS rn
        FROM ewma
        WHERE recorded_at >= CURRENT_DATE - INTERVAL '30 days'
    ) t
    WHERE rn = 1
) e
  ON p.id = e.product_id
ORDER BY e.ewma DESC
LIMIT 10;