-- +goose up
CREATE VIEW productdetail as
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

-- +goose down
DROP VIEW productdetail;