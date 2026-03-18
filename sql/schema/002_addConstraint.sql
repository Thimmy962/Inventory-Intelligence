-- +goose up
ALTER TABLE inventory
ALTER COLUMN quantity_on_hand SET NOT NULL;
