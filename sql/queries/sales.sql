-- name: CreateSales :one
INSERT INTO sales (total_amount, sale_date)
VALUES ($1, NOW()) RETURNING id;

-- name: CreateSalesItems :exec
INSERT INTO sales_items (sales_id, product_id, quantity_sold, price_at_sale)
VALUES ($1, $2, $3, $4);

-- name: DeleteSalesItems :exec
DELETE FROM sales_items WHERE sales_id = $1;

-- name: GetProductInventory :one
SELECT p.id, p.product_name, p.price, i.quantity_on_hand 
FROM products p
JOIN inventory i ON p.id = i.product_id
WHERE p.id = $1;

-- name: DeleteSales :exec
DELETE FROM sales WHERE id = $1;


-- name: CreateAdjustment :one
INSERT INTO adjustments (product_id, quantity_changed, reason, adjustment_date)
VALUES ($1, $2, $3, NOW()) RETURNING id;

-- name: DeleteAdjustment :exec
DELETE FROM adjustments WHERE id = $1;