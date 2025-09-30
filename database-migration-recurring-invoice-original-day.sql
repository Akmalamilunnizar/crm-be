-- Migration for Recurring Invoice Original Day
-- This migration adds the original_day column to preserve the template day for recurring invoices

-- Add original_day column to recurring_invoices table
ALTER TABLE `recurring_invoices` ADD COLUMN `original_day` int NOT NULL DEFAULT 1 COMMENT 'Preserve the original template day (1-31) for recurring invoice calculations';

-- Update existing records to set original_day based on their current invoice_date
UPDATE `recurring_invoices` SET `original_day` = DAY(`invoice_date`);

-- Add index for better performance on original_day queries
CREATE INDEX `idx_recurring_invoices_original_day` ON `recurring_invoices` (`original_day`);
