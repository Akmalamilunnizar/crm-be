-- Add frequency index
CREATE INDEX `idx_recurring_invoices_frequency` ON `recurring_invoices` (`frequency`);
