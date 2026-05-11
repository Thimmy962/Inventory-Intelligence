-- +goose up
ALTER TABLE sales
ALTER column staff_id SET NOT NULL;
