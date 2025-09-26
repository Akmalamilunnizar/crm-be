-- Add indexes for better performance
CREATE INDEX `idx_recurring_invoices_customer_status` ON `recurring_invoices` (`customer_id`, `status`);
