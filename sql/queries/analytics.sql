-- name: GetRecentProductEWMA :one
SELECT * FROM ewma WHERE product_id = $1
ORDER BY recorded_at DESC LIMIT 1;

-- name: AddNewEWMA :exec
INSERT INTO ewma (product_id, ewma, recorded_at)
VALUES ($1, $2, $3);