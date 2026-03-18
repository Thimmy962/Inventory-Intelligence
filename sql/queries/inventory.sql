-- name: NewInventory :exec
INSERT INTO inventory (product_id, quantity_on_hand, last_updated)
VALUES ($1, $2, NOW());

-- name: UpdatedInventory :exec
UPDATE inventory SET quantity_on_hand = $2, last_updated = NOW()
WHERE product_id = $1;


-- name: GetInventory :one
SELECT * FROM inventory WHERE product_id = $1;