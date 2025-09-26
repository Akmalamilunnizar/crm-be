-- Add due_date column to invoices table
ALTER TABLE `invoices` ADD COLUMN `due_date` DATE NULL AFTER `status`;

-- Add index for better performance on due_date queries
CREATE INDEX `idx_invoices_due_date` ON `invoices` (`due_date`);
