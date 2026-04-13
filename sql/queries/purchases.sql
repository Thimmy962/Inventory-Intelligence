-- name: CreatePurchase :one
INSERT INTO purchases (id, product_id, quantity_added, purchase_date)
VALUES (gen_random_uuid(), $1, $2, NOW())
Returning id;

-- name: DeletePurchase :exec
DELETE FROM purchases WHERE id = $1;


-- name: CreateAdjustments :one
INSERT INTO adjustments (product_id, quantity_changed, reason, adjustment_date)
VALUES ($1, $2, $3, NOW())
Returning product_id, quantity_changed;