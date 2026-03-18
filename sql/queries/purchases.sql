-- name: CreatePurchase :exec
INSERT INTO purchases (id, product_id, quantity_added, purchase_date)
VALUES (gen_random_uuid(), $1, $2, NOW());