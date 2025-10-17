-- Add deleted_at column to customers table for soft delete functionality
-- This allows preserving customer data and their related records (tickets, invoices, etc.)
-- when customers stop subscribing, while hiding them from normal queries

ALTER TABLE customers 
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

-- Add index for better performance on soft delete queries
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

-- Add comment to document the purpose
ALTER TABLE customers 
MODIFY COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL 
COMMENT 'Soft delete timestamp - NULL means active customer, timestamp means deleted customer';
