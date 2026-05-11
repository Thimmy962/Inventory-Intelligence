-- +goose up
-- 1. Add the column first
ALTER TABLE sales ADD COLUMN staff_id TEXT;

-- 2. Add the foreign key constraint
ALTER TABLE sales 
ADD CONSTRAINT fk_sales_staff 
FOREIGN KEY (staff_id) REFERENCES staffs(id);
