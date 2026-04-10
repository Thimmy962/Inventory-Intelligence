-- name: SearchProductForCheckout :many
SELECT p.id, p.product_name, p.price, i.quantity_on_hand
FROM products p JOIN inventory i ON p.id = i.product_id
WHERE to_tsvector('english', p.product_name) @@ websearch_to_tsquery('english', $1)
OR p.product_name ILIKE '%' || $1 || '%';

-- the function above aids searching the db for somethings