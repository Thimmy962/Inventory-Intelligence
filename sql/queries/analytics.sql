-- name: GetRecentProductEWMA :one
SELECT * FROM ewma WHERE product_id = $1
ORDER BY recorded_at DESC LIMIT 1;

-- name: AddNewEWMA :exec
INSERT INTO ewma (product_id, ewma, recorded_at)
VALUES ($1, $2, $3);

-- name: ProductTrend :many
SELECT recorded_at, ewma
FROM ewma
WHERE product_id = $1 AND recorded_at >= CURRENT_DATE - INTERVAL '30 days';